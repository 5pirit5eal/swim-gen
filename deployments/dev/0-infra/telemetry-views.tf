# =============================================================================
# Analytics Views for Looker Studio — Dev
# =============================================================================
#
# _AllLogs schema (linked log bucket dataset — LogEntry proto):
#
#   timestamp                          TIMESTAMP
#   log_id                             STRING
#   resource                           STRUCT<type STRING, labels JSON>
#     resource.labels                  JSON  → JSON_VALUE(resource.labels, '$.service_name')
#   json_payload                       JSON  → JSON_VALUE(json_payload, '$.user_id')
#   http_request                       STRUCT
#     .request_method                  STRING
#     .request_url                     STRING  (full URL incl. origin + query)
#     .status                          INT64
#     .latency                         STRUCT<seconds INT64, nanos INT64>
#
# Route extraction from request_url:
#   1. Extract path:  REGEXP_EXTRACT(request_url, r'https?://[^/]+(/[^?]*)')
#   2. Strip /api prefix for swim-gen-frontend and swim-gen-bff (routing
#      convention — the frontend calls /api/<path>, the BFF strips /api before
#      forwarding to the backend, so all three services should show the same
#      canonical route).  Backend routes are already clean.
#
# IMPORTANT: Views return no rows until traffic arrives. This is expected.
# =============================================================================

locals {
  linked_logs = "`${var.project_id}.swim_gen_dev_logs._AllLogs`"
}

# -----------------------------------------------------------------------------
# Monthly Active Users (per month, per year, per service)
# -----------------------------------------------------------------------------

resource "google_bigquery_table" "v_monthly_active_users" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_monthly_active_users"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      SELECT
        FORMAT_DATE('%Y-%m', DATE(timestamp))                  AS month,
        EXTRACT(YEAR  FROM timestamp)                          AS year,
        EXTRACT(MONTH FROM timestamp)                          AS month_num,
        JSON_VALUE(resource.labels, '$.service_name')          AS service,
        COUNT(DISTINCT JSON_VALUE(json_payload, '$.user_id'))  AS active_users
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/stdout"
        AND JSON_VALUE(json_payload, '$.user_id') IS NOT NULL
        -- Only count successful responses to prevent inflating MAU with failed/unauthorized requests
        AND CAST(JSON_VALUE(json_payload, '$.httpResponse.status') AS INT64) < 400
        AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
      GROUP BY month, year, month_num, service
      ORDER BY month DESC, service
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Total Unique Users (all-time, per year)
# -----------------------------------------------------------------------------

resource "google_bigquery_table" "v_total_users" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_total_users"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      -- All-time total
      SELECT
        'all_time'                                             AS period,
        NULL                                                   AS year_date,
        COUNT(DISTINCT JSON_VALUE(json_payload, '$.user_id')) AS unique_users,
        MIN(DATE(timestamp))                                   AS first_seen,
        MAX(DATE(timestamp))                                   AS last_seen
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/stdout"
        AND JSON_VALUE(json_payload, '$.user_id') IS NOT NULL
        -- Only count successful responses to prevent inflating counts with failed/unauthorized requests
        AND CAST(JSON_VALUE(json_payload, '$.httpResponse.status') AS INT64) < 400

      UNION ALL

      -- Per-year totals
      SELECT
        'yearly'                                                         AS period,
        DATE(EXTRACT(YEAR FROM timestamp), 1, 1)                        AS year_date,
        COUNT(DISTINCT JSON_VALUE(json_payload, '$.user_id'))           AS unique_users,
        MIN(DATE(timestamp))                                             AS first_seen,
        MAX(DATE(timestamp))                                             AS last_seen
      FROM ${local.linked_logs}
      WHERE
        log_id = "run.googleapis.com/stdout"
        AND JSON_VALUE(json_payload, '$.user_id') IS NOT NULL
        -- Only count successful responses to prevent inflating counts with failed/unauthorized requests
        AND CAST(JSON_VALUE(json_payload, '$.httpResponse.status') AS INT64) < 400
      GROUP BY year_date
      ORDER BY period, year_date DESC
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Request Volume by Day, Service, Route, and Status Class
# -----------------------------------------------------------------------------
# route = URL path without query string, extracted from http_request.request_url.
# Also includes a rolled-up total row per day+service (route = '_total').

resource "google_bigquery_table" "v_request_volume" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_request_volume"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      WITH base AS (
        SELECT
          DATE(timestamp)                                AS day,
          JSON_VALUE(resource.labels, '$.service_name') AS service,
          -- Strip /api prefix for frontend + bff; backend routes are already clean.
          IFNULL(
            CASE
              WHEN JSON_VALUE(resource.labels, '$.service_name')
                     IN ('swim-gen-frontend', 'swim-gen-bff')
              THEN REGEXP_REPLACE(
                     IFNULL(
                       REGEXP_EXTRACT(http_request.request_url, r'https?://[^/]+(/[^?]*)'),
                       '/'
                     ),
                     r'^/api', ''
                   )
              ELSE REGEXP_EXTRACT(http_request.request_url, r'https?://[^/]+(/[^?]*)')
            END,
            '/'
          )                                              AS route,
          http_request.status                            AS status_code,
          CONCAT(
            CAST(CAST(http_request.status / 100 AS INT64) AS STRING), 'xx'
          )                                              AS status_class
        FROM ${local.linked_logs}
        WHERE
          log_id = "run.googleapis.com/requests"
          AND http_request.status IS NOT NULL
          AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
          -- Exclude static asset requests (frontend noise: /_next/, /assets/, /public/, /wordpress/, *.html, *.svg, *.xml, *.php)
          AND NOT REGEXP_CONTAINS(
            http_request.request_url,
            r'(/_next/|/assets/|/public/|/wordpress/|\.(html|svg|xml|php)(\?|$))'
          )
      )

      -- Per-route breakdown
      SELECT day, service, route, status_code, status_class, COUNT(*) AS request_count
      FROM base
      GROUP BY day, service, route, status_code, status_class

      UNION ALL

      -- Daily total per service (all routes combined)
      SELECT day, service, '_total' AS route, status_code, status_class, COUNT(*) AS request_count
      FROM base
      GROUP BY day, service, status_code, status_class

      ORDER BY day DESC, service, route, status_code
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Request Latency Percentiles by Day, Service, and Route
# -----------------------------------------------------------------------------
# http_request.latency is STRUCT<seconds INT64, nanos INT64>.
# Also includes a rolled-up _total row per day+service.

resource "google_bigquery_table" "v_request_latency" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_request_latency"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      WITH base AS (
        SELECT
          DATE(timestamp)                                AS day,
          JSON_VALUE(resource.labels, '$.service_name') AS service,
          -- Strip /api prefix for frontend + bff; backend routes are already clean.
          IFNULL(
            CASE
              WHEN JSON_VALUE(resource.labels, '$.service_name')
                     IN ('swim-gen-frontend', 'swim-gen-bff')
              THEN REGEXP_REPLACE(
                     IFNULL(
                       REGEXP_EXTRACT(http_request.request_url, r'https?://[^/]+(/[^?]*)'),
                       '/'
                     ),
                     r'^/api', ''
                   )
              ELSE REGEXP_EXTRACT(http_request.request_url, r'https?://[^/]+(/[^?]*)')
            END,
            '/'
          )                                              AS route,
          http_request.latency.seconds * 1000.0
            + http_request.latency.nanos / 1000000.0    AS latency_ms
        FROM ${local.linked_logs}
        WHERE
          log_id = "run.googleapis.com/requests"
          AND http_request.latency IS NOT NULL
          AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
          -- Exclude static asset requests (frontend noise: /_next/, /assets/, /public/, /wordpress/, *.html, *.svg, *.xml, *.php)
          AND NOT REGEXP_CONTAINS(
            http_request.request_url,
            r'(/_next/|/assets/|/public/|/wordpress/|\.(html|svg|xml|php)(\?|$))'
          )
      )

      -- Per-route latency
      SELECT
        day, service, route,
        COUNT(*)                                                            AS request_count,
        ROUND(APPROX_QUANTILES(latency_ms, 100)[OFFSET(50)], 2)           AS p50_ms,
        ROUND(APPROX_QUANTILES(latency_ms, 100)[OFFSET(95)], 2)           AS p95_ms,
        ROUND(APPROX_QUANTILES(latency_ms, 100)[OFFSET(99)], 2)           AS p99_ms
      FROM base
      GROUP BY day, service, route

      UNION ALL

      -- Daily total per service (all routes combined)
      SELECT
        day, service, '_total' AS route,
        COUNT(*)                                                            AS request_count,
        ROUND(APPROX_QUANTILES(latency_ms, 100)[OFFSET(50)], 2)           AS p50_ms,
        ROUND(APPROX_QUANTILES(latency_ms, 100)[OFFSET(95)], 2)           AS p95_ms,
        ROUND(APPROX_QUANTILES(latency_ms, 100)[OFFSET(99)], 2)           AS p99_ms
      FROM base
      GROUP BY day, service

      ORDER BY day DESC, service, route
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}

# -----------------------------------------------------------------------------
# Error Rate by Day, Service, and Route
# -----------------------------------------------------------------------------
# Also includes a rolled-up _total row per day+service.

resource "google_bigquery_table" "v_error_rate" {
  project    = var.project_id
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "v_error_rate"

  deletion_protection = false

  view {
    use_legacy_sql = false
    query          = <<-SQL
      WITH base AS (
        SELECT
          DATE(timestamp)                                AS day,
          JSON_VALUE(resource.labels, '$.service_name') AS service,
          -- Strip /api prefix for frontend + bff; backend routes are already clean.
          IFNULL(
            CASE
              WHEN JSON_VALUE(resource.labels, '$.service_name')
                     IN ('swim-gen-frontend', 'swim-gen-bff')
              THEN REGEXP_REPLACE(
                     IFNULL(
                       REGEXP_EXTRACT(http_request.request_url, r'https?://[^/]+(/[^?]*)'),
                       '/'
                     ),
                     r'^/api', ''
                   )
              ELSE REGEXP_EXTRACT(http_request.request_url, r'https?://[^/]+(/[^?]*)')
            END,
            '/'
          )                                              AS route,
          http_request.status                            AS status
        FROM ${local.linked_logs}
        WHERE
          log_id = "run.googleapis.com/requests"
          AND http_request.status IS NOT NULL
          AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1825 DAY)
          -- Exclude static asset requests (frontend noise: /_next/, /assets/, /public/, /wordpress/, *.html, *.svg, *.xml, *.php)
          AND NOT REGEXP_CONTAINS(
            http_request.request_url,
            r'(/_next/|/assets/|/public/|/wordpress/|\.(html|svg|xml|php)(\?|$))'
          )
      )

      -- Per-route error rate
      SELECT
        day, service, route,
        COUNT(*)                                                        AS total_requests,
        COUNTIF(status >= 500)                                          AS server_errors,
        COUNTIF(status >= 400 AND status < 500)                        AS client_errors,
        ROUND(SAFE_DIVIDE(COUNTIF(status >= 500), COUNT(*)) * 100, 4)  AS server_error_rate_pct,
        ROUND(SAFE_DIVIDE(COUNTIF(status >= 400), COUNT(*)) * 100, 4)  AS error_rate_pct
      FROM base
      GROUP BY day, service, route

      UNION ALL

      -- Daily total per service (all routes combined)
      SELECT
        day, service, '_total' AS route,
        COUNT(*)                                                        AS total_requests,
        COUNTIF(status >= 500)                                          AS server_errors,
        COUNTIF(status >= 400 AND status < 500)                        AS client_errors,
        ROUND(SAFE_DIVIDE(COUNTIF(status >= 500), COUNT(*)) * 100, 4)  AS server_error_rate_pct,
        ROUND(SAFE_DIVIDE(COUNTIF(status >= 400), COUNT(*)) * 100, 4)  AS error_rate_pct
      FROM base
      GROUP BY day, service

      ORDER BY day DESC, service, route
    SQL
  }

  depends_on = [google_bigquery_dataset.analytics, google_logging_linked_dataset.telemetry]
}
