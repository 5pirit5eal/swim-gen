# --------------------------------------------------------------------------
# BigQuery Views for Looker Studio Dashboard
# --------------------------------------------------------------------------
# These views query tables auto-created by the Cloud Logging sink.
# Table names are determined by Cloud Logging:
#   - run_googleapis_com_requests  (Cloud Run auto-generated request logs)
#   - run_googleapis_com_stdout    (application stdout logs)
#
# Column naming: BigQuery Cloud Logging export uses snake_case columns
#   - json_payload (not jsonPayload)
#   - http_request (not httpRequest)
#   - resource.labels.service_name
#
# IMPORTANT: These views will fail until the logging sink has created
# the underlying tables (i.e., after first deployment + traffic).
# --------------------------------------------------------------------------

# --------------------------------------------------------------------------
# Active Users per Day
# --------------------------------------------------------------------------
# Uses application stdout logs where user_id is present.
# Both BFF and backend emit user_id in structured logs.

resource "google_bigquery_table" "v_active_users" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.telemetry.dataset_id
  table_id   = "v_active_users"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp) AS day,
        COUNT(DISTINCT JSON_VALUE(json_payload, '$.user_id')) AS daily_active_users
      FROM `${var.project_id}.${google_bigquery_dataset.telemetry.dataset_id}.run_googleapis_com_stdout`
      WHERE
        JSON_VALUE(json_payload, '$.user_id') IS NOT NULL
        AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 365 DAY)
      GROUP BY day
      ORDER BY day
    SQL
  }

  depends_on = [google_bigquery_dataset.telemetry]
}

# --------------------------------------------------------------------------
# Request Latency Percentiles per Day and Service
# --------------------------------------------------------------------------
# Uses Cloud Run auto-generated request logs.
# http_request.latency is a duration string (e.g. "0.115239727s").

resource "google_bigquery_table" "v_request_latency" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.telemetry.dataset_id
  table_id   = "v_request_latency"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp) AS day,
        resource.labels.service_name AS service,
        APPROX_QUANTILES(
          CAST(REPLACE(http_request.latency, 's', '') AS FLOAT64) * 1000, 100
        )[OFFSET(50)] AS p50_ms,
        APPROX_QUANTILES(
          CAST(REPLACE(http_request.latency, 's', '') AS FLOAT64) * 1000, 100
        )[OFFSET(95)] AS p95_ms,
        APPROX_QUANTILES(
          CAST(REPLACE(http_request.latency, 's', '') AS FLOAT64) * 1000, 100
        )[OFFSET(99)] AS p99_ms
      FROM `${var.project_id}.${google_bigquery_dataset.telemetry.dataset_id}.run_googleapis_com_requests`
      WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 365 DAY)
      GROUP BY day, service
      ORDER BY day
    SQL
  }

  depends_on = [google_bigquery_dataset.telemetry]
}

# --------------------------------------------------------------------------
# Request Volume per Day, Service, and Status Code
# --------------------------------------------------------------------------

resource "google_bigquery_table" "v_request_volume" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.telemetry.dataset_id
  table_id   = "v_request_volume"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp) AS day,
        resource.labels.service_name AS service,
        http_request.status AS status_code,
        COUNT(*) AS request_count
      FROM `${var.project_id}.${google_bigquery_dataset.telemetry.dataset_id}.run_googleapis_com_requests`
      WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 365 DAY)
      GROUP BY day, service, status_code
      ORDER BY day
    SQL
  }

  depends_on = [google_bigquery_dataset.telemetry]
}
