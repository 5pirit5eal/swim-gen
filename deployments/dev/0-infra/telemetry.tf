# =============================================================================
# Telemetry Infrastructure — Dev
# =============================================================================
#
# Architecture:
#
#   Cloud Run logs
#     └─► Log Sink ──► Custom Log Bucket (analytics-enabled, 5yr retention)
#                          ├─ Observability Analytics  (SQL in GCP console)
#                          └─ Linked BigQuery Dataset  (for Looker Studio)
#                               └─ _AllLogs virtual view
#
#   Custom views for Looker Studio live in a separate google_bigquery_dataset
#   ("analytics") because the linked dataset is read-only.
#
# =============================================================================

# -----------------------------------------------------------------------------
# 1. Custom Log Bucket — analytics-enabled, 5-year retention
# -----------------------------------------------------------------------------
# Routes only application-relevant logs (stdout/stderr + Cloud Run request
# logs). System/lifecycle events are excluded to reduce noise and storage.
#
# NOTE: enable_analytics cannot be disabled once set. Location cannot be
# changed after creation. Both are intentional.

resource "google_logging_project_bucket_config" "cloud_run" {
  project          = var.project_id
  location         = var.region
  bucket_id        = "swim-gen-cloud-run"
  description      = "Cloud Run application logs (stdout, stderr, requests) — 5yr retention"
  retention_days   = 1825 # 5 years
  enable_analytics = true

  depends_on = [google_project_service.apis]
}

# -----------------------------------------------------------------------------
# 2. Log Sink → Custom Log Bucket
# -----------------------------------------------------------------------------
# Captures Cloud Run stdout/stderr (structured app logs) and the
# auto-generated request logs. System/lifecycle logs are excluded.

resource "google_logging_project_sink" "cloud_run_logs" {
  name    = "cloud-run-logs-to-bucket"
  project = var.project_id

  destination = "logging.googleapis.com/${google_logging_project_bucket_config.cloud_run.id}"

  # Frontend (swim-gen-frontend) is excluded from stdout/stderr — it produces
  # high-volume SSR logs that are not useful for long-term analysis.
  # Request logs from the frontend ARE still captured for traffic/error tracking.
  filter = <<-FILTER
    resource.type="cloud_run_revision"
    AND (
      (
        log_id("run.googleapis.com/stdout")
        OR log_id("run.googleapis.com/stderr")
      )
      AND resource.labels.service_name!="swim-gen-frontend"
    )
    OR log_id("run.googleapis.com/requests")
  FILTER

  # unique_writer_identity is not set for same-project log bucket destinations.
  # Cloud Logging automatically uses its own service account; no IAM binding needed.

  depends_on = [google_logging_project_bucket_config.cloud_run]
}

# -----------------------------------------------------------------------------
# 3. Linked BigQuery Dataset on the Log Bucket
# -----------------------------------------------------------------------------
# Creates a read-only BigQuery dataset backed by the log bucket.
# Cloud Logging automatically maps log views → BigQuery virtual views:
#   _AllLogs  →  PROJECT.swim_gen_dev_logs._AllLogs
#
# This dataset can be queried from Looker Studio (BigQuery connector) and
# from Observability Analytics in the GCP console.

resource "google_logging_linked_dataset" "telemetry" {
  link_id     = "swim_gen_dev_logs"
  bucket      = google_logging_project_bucket_config.cloud_run.id
  description = "Linked BigQuery dataset for Cloud Run log analytics (dev)"

  depends_on = [google_logging_project_bucket_config.cloud_run]
}

# -----------------------------------------------------------------------------
# 4. Analytics BigQuery Dataset — custom views for Looker Studio
# -----------------------------------------------------------------------------
# The linked dataset is read-only; we cannot create views inside it.
# This writable dataset holds our pre-built views that reference _AllLogs.

resource "google_bigquery_dataset" "analytics" {
  dataset_id  = "swim_gen_dev_analytics"
  description = "Pre-built analytics views for Looker Studio (sourced from log bucket linked dataset)"
  project     = var.project_id
  location    = var.region

  delete_contents_on_destroy = true # dev only — safe to destroy

  depends_on = [google_project_service.apis]
}

# =============================================================================
# _Trace Linked BigQuery Dataset (Cloud Trace / Observability Analytics)
# =============================================================================
# Cloud Trace stores spans in a system-managed observability bucket (_Trace).
# The _AllSpans view is auto-created when the first trace arrives.
#
# Terraform does NOT yet support google_observability_linked_dataset.
# To expose traces in BigQuery / Looker Studio, run once per project:
#
#   gcloud beta observability buckets datasets links create swim_gen_dev_traces \
#     --dataset=Spans \
#     --bucket=_Trace \
#     --location=global \
#     --project=rubenschulze-sandbox
#
# After creation, reference the view as:
#   `rubenschulze-sandbox.swim_gen_dev_traces._AllSpans`
# =============================================================================
