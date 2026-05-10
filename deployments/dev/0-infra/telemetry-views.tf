# =============================================================================
# Analytics Views for Looker Studio — Dev
# =============================================================================
#
# These views live in the writable `swim_gen_dev_analytics` dataset and query
# the read-only linked dataset `swim_gen_dev_logs` (backed by the log bucket).
#
# _AllLogs schema (linked log bucket dataset — LogEntry proto representation):
#
#   timestamp            TIMESTAMP
#   log_id               STRING
#   severity             STRING
#   resource             JSON        ← nested, use JSON_VALUE()
#   json_payload         JSON        ← nested, use JSON_VALUE()
#   http_request         STRUCT      ← direct field access
#     .request_method    STRING
#     .request_url       STRING
#     .status            INT64
#     .latency           STRUCT<seconds INT64, nanos INT64>  ← Duration proto
#     .user_agent        STRING
#     .remote_ip         STRING
#
# Extract service name:  JSON_VALUE(resource, '$.labels.service_name')
# Extract latency ms:    http_request.latency.seconds * 1000
#                        + http_request.latency.nanos / 1000000
# Extract user_id:       JSON_VALUE(json_payload, '$.user_id')
#
# IMPORTANT: Views return no rows until traffic arrives. This is expected.
# =============================================================================

locals {
  linked_logs = "`${var.project_id}.swim_gen_dev_logs._AllLogs`"
}

# -----------------------------------------------------------------------------
# Daily Active Users
# -----------------------------------------------------------------------------
# Counts distinct authenticated users per day from structured app logs.
# Both BFF and Backend emit user_id (JWT sub) in structured stdout logs.

resource "google_bigquery_table" "v_daily_active_users" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_daily_active_users"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)                                        AS day,
        JSON_VALUE(resource, '$.labels.service_name')         AS service,
        COUNT(DISTINCT JSON_VALUE(json_payload, '$.user_id')) AS active_users
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/stdout"
        AND JSON_VALUE(json_payload, '$.user_id') IS NOT NULL
        AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
      GROUP BY day, service
      ORDER BY day DESC
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Total (All-Time) Unique Users
# -----------------------------------------------------------------------------
# Every distinct user_id ever seen — proxy for total registered active users.

resource "google_bigquery_table" "v_total_users" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_total_users"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        COUNT(DISTINCT JSON_VALUE(json_payload, '$.user_id')) AS total_unique_users,
        MIN(DATE(timestamp))                                  AS first_seen,
        MAX(DATE(timestamp))                                  AS last_seen
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/stdout"
        AND JSON_VALUE(json_payload, '$.user_id') IS NOT NULL
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Request Volume by Day, Service, and Status Class
# -----------------------------------------------------------------------------
# Cloud Run request logs populate http_request as a STRUCT.
# http_request.status is INT64.

resource "google_bigquery_table" "v_request_volume" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_request_volume"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)                               AS day,
        JSON_VALUE(resource, '$.labels.service_name') AS service,
        http_request.status                           AS status_code,
        CONCAT(
          CAST(CAST(http_request.status / 100 AS INT64) AS STRING), 'xx'
        )                                             AS status_class,
        COUNT(*)                                      AS request_count
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/requests"
        AND http_request.status IS NOT NULL
        AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
      GROUP BY day, service, status_code, status_class
      ORDER BY day DESC, service, status_code
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Request Latency Percentiles by Day and Service
# -----------------------------------------------------------------------------
# http_request.latency is STRUCT<seconds INT64, nanos INT64> (Duration proto).
# Convert to milliseconds: seconds * 1000 + nanos / 1_000_000

resource "google_bigquery_table" "v_request_latency" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_request_latency"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)                               AS day,
        JSON_VALUE(resource, '$.labels.service_name') AS service,
        COUNT(*)                                      AS request_count,
        ROUND(
          APPROX_QUANTILES(
            http_request.latency.seconds * 1000.0
              + http_request.latency.nanos / 1000000.0,
            100
          )[OFFSET(50)], 2
        )                                             AS p50_ms,
        ROUND(
          APPROX_QUANTILES(
            http_request.latency.seconds * 1000.0
              + http_request.latency.nanos / 1000000.0,
            100
          )[OFFSET(95)], 2
        )                                             AS p95_ms,
        ROUND(
          APPROX_QUANTILES(
            http_request.latency.seconds * 1000.0
              + http_request.latency.nanos / 1000000.0,
            100
          )[OFFSET(99)], 2
        )                                             AS p99_ms
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/requests"
        AND http_request.latency IS NOT NULL
        AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
      GROUP BY day, service
      ORDER BY day DESC, service
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Error Rate by Day and Service
# -----------------------------------------------------------------------------

resource "google_bigquery_table" "v_error_rate" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_error_rate"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)                               AS day,
        JSON_VALUE(resource, '$.labels.service_name') AS service,
        COUNT(*)                                      AS total_requests,
        COUNTIF(http_request.status >= 500)           AS server_errors,
        COUNTIF(http_request.status >= 400
          AND http_request.status < 500)              AS client_errors,
        ROUND(
          SAFE_DIVIDE(
            COUNTIF(http_request.status >= 500),
            COUNT(*)
          ) * 100, 4
        )                                             AS server_error_rate_pct,
        ROUND(
          SAFE_DIVIDE(
            COUNTIF(http_request.status >= 400),
            COUNT(*)
          ) * 100, 4
        )                                             AS error_rate_pct
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/requests"
        AND http_request.status IS NOT NULL
        AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
      GROUP BY day, service
      ORDER BY day DESC, service
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}
