# =============================================================================
# Telemetry Infrastructure — Prod
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

resource "google_logging_project_sink" "cloud_run_logs" {
  name    = "cloud-run-logs-to-bucket"
  project = var.project_id

  destination = "logging.googleapis.com/${google_logging_project_bucket_config.cloud_run.id}"

  filter = <<-FILTER
    resource.type="cloud_run_revision"
    AND (
      log_id("run.googleapis.com/stdout")
      OR log_id("run.googleapis.com/stderr")
      OR log_id("run.googleapis.com/requests")
    )
  FILTER

  unique_writer_identity = true

  depends_on = [google_logging_project_bucket_config.cloud_run]
}

resource "google_project_iam_member" "log_sink_bucket_writer" {
  project = var.project_id
  role    = "roles/logging.bucketWriter"
  member  = google_logging_project_sink.cloud_run_logs.writer_identity
}

# -----------------------------------------------------------------------------
# 3. Linked BigQuery Dataset on the Log Bucket
# -----------------------------------------------------------------------------

resource "google_logging_linked_dataset" "telemetry" {
  link_id     = "swim_gen_prod_logs"
  bucket      = google_logging_project_bucket_config.cloud_run.id
  description = "Linked BigQuery dataset for Cloud Run log analytics (prod)"

  depends_on = [google_logging_project_bucket_config.cloud_run]
}

# -----------------------------------------------------------------------------
# 4. Analytics BigQuery Dataset — custom views for Looker Studio
# -----------------------------------------------------------------------------

resource "google_bigquery_dataset" "analytics" {
  dataset_id  = "swim_gen_prod_analytics"
  description = "Pre-built analytics views for Looker Studio (sourced from log bucket linked dataset)"
  project     = var.project_id
  location    = var.region

  delete_contents_on_destroy = false # prod — protect data on destroy

  depends_on = [google_project_service.apis]
}

# =============================================================================
# _Trace Linked BigQuery Dataset (Cloud Trace / Observability Analytics)
# =============================================================================
# Terraform does NOT yet support google_observability_linked_dataset.
# Run once per project after initial deployment:
#
#   gcloud beta observability buckets datasets links create swim_gen_prod_traces \
#     --dataset=Spans \
#     --bucket=_Trace \
#     --location=global \
#     --project=swim-gen-prod
#
# After creation, reference the view as:
#   `swim-gen-prod.swim_gen_prod_traces._AllSpans`
# =============================================================================
