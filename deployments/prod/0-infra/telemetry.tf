# --------------------------------------------------------------------------
# BigQuery Dataset for Telemetry (Cloud Logging sink destination)
# --------------------------------------------------------------------------

resource "google_bigquery_dataset" "telemetry" {
  dataset_id  = "swim_gen_prod_telemetry"
  description = "Telemetry data from Cloud Run services (logs, metrics)"
  project     = var.project_id
  location    = var.region

  # Keep data indefinitely for 1-year+ analysis
  # default_table_expiration_ms not set

  delete_contents_on_destroy = false # prod — protect data on destroy

  depends_on = [google_project_service.apis]
}

# --------------------------------------------------------------------------
# Cloud Logging Sink → BigQuery
# --------------------------------------------------------------------------

resource "google_logging_project_sink" "cloud_run_logs" {
  name        = "cloud-run-logs-to-bigquery"
  project     = var.project_id
  destination = "bigquery.googleapis.com/projects/${var.project_id}/datasets/${google_bigquery_dataset.telemetry.dataset_id}"
  filter      = "resource.type=\"cloud_run_revision\""

  unique_writer_identity = true

  bigquery_options {
    use_partitioned_tables = true
  }

  depends_on = [google_project_service.apis]
}

# --------------------------------------------------------------------------
# Grant the log sink's writer identity BigQuery Data Editor on the dataset
# --------------------------------------------------------------------------

resource "google_bigquery_dataset_iam_member" "log_sink_writer" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.telemetry.dataset_id
  role       = "roles/bigquery.dataEditor"
  member     = google_logging_project_sink.cloud_run_logs.writer_identity
}

# --------------------------------------------------------------------------
# Cloud Monitoring Metrics Export to BigQuery
# --------------------------------------------------------------------------
# Terraform does not currently support configuring Cloud Monitoring
# metrics export to BigQuery natively. To retain metrics beyond the
# default 6-week window:
#
# 1. Go to GCP Console → Monitoring → Settings → Metrics management
# 2. Configure a Metrics Export to BigQuery targeting the telemetry dataset
#
# Alternatively, use gcloud:
#   gcloud alpha monitoring metrics-scopes create ...
#
# This is a one-time manual configuration per project.
# --------------------------------------------------------------------------
