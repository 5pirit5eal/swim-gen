# =============================================================================
# Analytics Views for Looker Studio — Dev
# =============================================================================
#
# These views live in the writable `swim_gen_dev_analytics` dataset and query
# the read-only linked dataset `swim_gen_dev_logs` (backed by the log bucket).
#
# Schema reference — LogEntry fields available in _AllLogs:
#   timestamp                          TIMESTAMP
#   log_id                             STRING  ("run.googleapis.com/stdout" etc.)
#   json_payload                       JSON    (structured app logs)
#   http_request.request_method        STRING
#   http_request.request_url           STRING
#   http_request.status                INT64
#   http_request.latency               STRING  ("0.115239727s")
#   resource.labels.service_name       STRING  (Cloud Run service)
#   resource.labels.revision_name      STRING
#
# App-log user context (emitted by both BFF and Backend):
#   JSON_VALUE(json_payload, '$.user_id')                    STRING
#   JSON_VALUE(json_payload, '$.httpRequest.method')         STRING
#   JSON_VALUE(json_payload, '$.httpRequest.path')           STRING
#   JSON_VALUE(json_payload, '$.httpResponse.status')        STRING
#   JSON_VALUE(json_payload, '$.httpResponse.elapsed')       STRING (ms as number)
#
# Linked dataset fully-qualified view:
#   `rubenschulze-sandbox.swim_gen_dev_logs._AllLogs`
#
# IMPORTANT: Views will return no rows until traffic arrives and logs are
# routed to the bucket. This is expected — views are created at apply time
# but data arrives later.
# =============================================================================

locals {
  linked_logs = "`${var.project_id}.swim_gen_dev_logs._AllLogs`"
}

# -----------------------------------------------------------------------------
# Daily Active Users
# -----------------------------------------------------------------------------
# Counts distinct authenticated users per day from structured app logs.
# A user is "active" if they made at least one authenticated request.
# Both BFF and Backend emit user_id from the JWT sub claim.

resource "google_bigquery_table" "v_daily_active_users" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_daily_active_users"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)                                       AS day,
        resource.labels.service_name                         AS service,
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
# Counts every distinct user_id ever seen across all time. Useful as a
# "total registered active users" proxy (users who have logged in at least once).

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
# Uses Cloud Run auto-generated request logs (http_request struct is populated).
# Groups by HTTP status class (2xx, 4xx, 5xx) for error-rate visibility.

resource "google_bigquery_table" "v_request_volume" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_request_volume"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)               AS day,
        resource.labels.service_name  AS service,
        http_request.status           AS status_code,
        CONCAT(
          CAST(CAST(http_request.status / 100 AS INT64) AS STRING), 'xx'
        )                             AS status_class,
        COUNT(*)                      AS request_count
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
# Uses Cloud Run request logs. http_request.latency is a duration string
# like "0.115239727s" — stripped of "s" and converted to milliseconds.

resource "google_bigquery_table" "v_request_latency" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_request_latency"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)              AS day,
        resource.labels.service_name AS service,
        COUNT(*)                     AS request_count,
        ROUND(
          APPROX_QUANTILES(
            CAST(REGEXP_REPLACE(http_request.latency, r's$', '') AS FLOAT64) * 1000,
            100
          )[OFFSET(50)], 2
        )                            AS p50_ms,
        ROUND(
          APPROX_QUANTILES(
            CAST(REGEXP_REPLACE(http_request.latency, r's$', '') AS FLOAT64) * 1000,
            100
          )[OFFSET(95)], 2
        )                            AS p95_ms,
        ROUND(
          APPROX_QUANTILES(
            CAST(REGEXP_REPLACE(http_request.latency, r's$', '') AS FLOAT64) * 1000,
            100
          )[OFFSET(99)], 2
        )                            AS p99_ms
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
# Computes error rate (5xx / total) and client-error rate (4xx / total).
# Useful for SLO tracking in Looker Studio.

resource "google_bigquery_table" "v_error_rate" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_error_rate"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        DATE(timestamp)              AS day,
        resource.labels.service_name AS service,
        COUNT(*)                     AS total_requests,
        COUNTIF(http_request.status >= 500)             AS server_errors,
        COUNTIF(http_request.status >= 400
          AND http_request.status < 500)                AS client_errors,
        ROUND(
          SAFE_DIVIDE(
            COUNTIF(http_request.status >= 500),
            COUNT(*)
          ) * 100, 4
        )                            AS server_error_rate_pct,
        ROUND(
          SAFE_DIVIDE(
            COUNTIF(http_request.status >= 400),
            COUNT(*)
          ) * 100, 4
        )                            AS error_rate_pct
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
