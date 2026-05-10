# =============================================================================
# Cloud Monitoring Dashboard — Dev
# =============================================================================
#
# Panels:
#   Row 1 — Traffic:    Request Count (by response class) | Request Latency (p50/p95/p99)
#   Row 2 — Errors:     Error Rate (4xx + 5xx)            | Server Error Rate (5xx only)
#   Row 3 — Resources:  CPU Utilization (per service)     | Memory Utilization (per service)
#   Row 4 — Scaling:    Instance Count (per service)      | Startup Latency
#
# Metric types (all under run.googleapis.com/):
#   request_count           DELTA/INT64  — requests per minute by response_code_class
#   request_latencies       DELTA/DIST   — latency distribution in ms
#   container/cpu/utilizations    GAUGE/DOUBLE — CPU % (0–100)
#   container/memory/utilizations GAUGE/DOUBLE — Memory % (0–100)
#   container/instance_count      GAUGE/INT64  — running instances
#   container/startup_latency     DELTA/DIST   — cold-start latency in ms
#
# NOTE: Do not edit this dashboard in the Cloud Console — Terraform will
# detect drift. To modify it, update the JSON here and re-apply.
# =============================================================================

resource "google_monitoring_dashboard" "swim_gen" {
  project = var.project_id
  dashboard_json = jsonencode({
    displayName = "swim-gen — Cloud Run (dev)"
    mosaicLayout = {
      columns = 12
      tiles = [

        # ── Row 1, col 0-5: Request Count ──────────────────────────────────
        {
          xPos   = 0
          yPos   = 0
          width  = 6
          height = 4
          widget = {
            title = "Request Count (per minute)"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = join(" ", [
                      "metric.type=\"run.googleapis.com/request_count\"",
                      "resource.type=\"cloud_run_revision\"",
                    ])
                    aggregation = {
                      alignmentPeriod    = "60s"
                      perSeriesAligner   = "ALIGN_RATE"
                      crossSeriesReducer = "REDUCE_SUM"
                      groupByFields = [
                        "resource.labels.service_name",
                        "metric.labels.response_code_class",
                      ]
                    }
                  }
                }
                plotType       = "LINE"
                legendTemplate = "$${resource.labels.service_name} $${metric.labels.response_code_class}"
              }]
              yAxis = { label = "requests/s", scale = "LINEAR" }
            }
          }
        },

        # ── Row 1, col 6-11: Request Latency p50/p95/p99 ───────────────────
        {
          xPos   = 6
          yPos   = 0
          width  = 6
          height = 4
          widget = {
            title = "Request Latency — p50 / p95 / p99 (ms)"
            xyChart = {
              dataSets = [
                {
                  timeSeriesQuery = {
                    timeSeriesFilter = {
                      filter = join(" ", [
                        "metric.type=\"run.googleapis.com/request_latencies\"",
                        "resource.type=\"cloud_run_revision\"",
                      ])
                      aggregation = {
                        alignmentPeriod    = "60s"
                        perSeriesAligner   = "ALIGN_DELTA"
                        crossSeriesReducer = "REDUCE_PERCENTILE_50"
                        groupByFields      = ["resource.labels.service_name"]
                      }
                    }
                  }
                  plotType       = "LINE"
                  legendTemplate = "$${resource.labels.service_name} p50"
                },
                {
                  timeSeriesQuery = {
                    timeSeriesFilter = {
                      filter = join(" ", [
                        "metric.type=\"run.googleapis.com/request_latencies\"",
                        "resource.type=\"cloud_run_revision\"",
                      ])
                      aggregation = {
                        alignmentPeriod    = "60s"
                        perSeriesAligner   = "ALIGN_DELTA"
                        crossSeriesReducer = "REDUCE_PERCENTILE_95"
                        groupByFields      = ["resource.labels.service_name"]
                      }
                    }
                  }
                  plotType       = "LINE"
                  legendTemplate = "$${resource.labels.service_name} p95"
                },
                {
                  timeSeriesQuery = {
                    timeSeriesFilter = {
                      filter = join(" ", [
                        "metric.type=\"run.googleapis.com/request_latencies\"",
                        "resource.type=\"cloud_run_revision\"",
                      ])
                      aggregation = {
                        alignmentPeriod    = "60s"
                        perSeriesAligner   = "ALIGN_DELTA"
                        crossSeriesReducer = "REDUCE_PERCENTILE_99"
                        groupByFields      = ["resource.labels.service_name"]
                      }
                    }
                  }
                  plotType       = "LINE"
                  legendTemplate = "$${resource.labels.service_name} p99"
                },
              ]
              yAxis = { label = "ms", scale = "LINEAR" }
            }
          }
        },

        # ── Row 2, col 0-5: 4xx Rate ───────────────────────────────────────
        {
          xPos   = 0
          yPos   = 4
          width  = 6
          height = 4
          widget = {
            title = "Client Error Rate (4xx requests/s)"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = join(" ", [
                      "metric.type=\"run.googleapis.com/request_count\"",
                      "resource.type=\"cloud_run_revision\"",
                      "metric.labels.response_code_class=\"4xx\"",
                    ])
                    aggregation = {
                      alignmentPeriod    = "60s"
                      perSeriesAligner   = "ALIGN_RATE"
                      crossSeriesReducer = "REDUCE_SUM"
                      groupByFields      = ["resource.labels.service_name"]
                    }
                  }
                }
                plotType       = "LINE"
                legendTemplate = "$${resource.labels.service_name}"
              }]
              yAxis = { label = "4xx/s", scale = "LINEAR" }
            }
          }
        },

        # ── Row 2, col 6-11: 5xx Rate ──────────────────────────────────────
        {
          xPos   = 6
          yPos   = 4
          width  = 6
          height = 4
          widget = {
            title = "Server Error Rate (5xx requests/s)"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = join(" ", [
                      "metric.type=\"run.googleapis.com/request_count\"",
                      "resource.type=\"cloud_run_revision\"",
                      "metric.labels.response_code_class=\"5xx\"",
                    ])
                    aggregation = {
                      alignmentPeriod    = "60s"
                      perSeriesAligner   = "ALIGN_RATE"
                      crossSeriesReducer = "REDUCE_SUM"
                      groupByFields      = ["resource.labels.service_name"]
                    }
                  }
                }
                plotType       = "LINE"
                legendTemplate = "$${resource.labels.service_name}"
              }]
              yAxis = { label = "5xx/s", scale = "LINEAR" }
            }
          }
        },

        # ── Row 3, col 0-5: CPU Utilization ────────────────────────────────
        {
          xPos   = 0
          yPos   = 8
          width  = 6
          height = 4
          widget = {
            title = "Container CPU Utilization (%)"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = join(" ", [
                      "metric.type=\"run.googleapis.com/container/cpu/utilizations\"",
                      "resource.type=\"cloud_run_revision\"",
                    ])
                    aggregation = {
                      alignmentPeriod    = "60s"
                      perSeriesAligner   = "ALIGN_MEAN"
                      crossSeriesReducer = "REDUCE_MEAN"
                      groupByFields      = ["resource.labels.service_name"]
                    }
                  }
                }
                plotType       = "LINE"
                legendTemplate = "$${resource.labels.service_name}"
              }]
              yAxis = { label = "%", scale = "LINEAR" }
            }
          }
        },

        # ── Row 3, col 6-11: Memory Utilization ────────────────────────────
        {
          xPos   = 6
          yPos   = 8
          width  = 6
          height = 4
          widget = {
            title = "Container Memory Utilization (%)"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = join(" ", [
                      "metric.type=\"run.googleapis.com/container/memory/utilizations\"",
                      "resource.type=\"cloud_run_revision\"",
                    ])
                    aggregation = {
                      alignmentPeriod    = "60s"
                      perSeriesAligner   = "ALIGN_MEAN"
                      crossSeriesReducer = "REDUCE_MEAN"
                      groupByFields      = ["resource.labels.service_name"]
                    }
                  }
                }
                plotType       = "LINE"
                legendTemplate = "$${resource.labels.service_name}"
              }]
              yAxis = { label = "%", scale = "LINEAR" }
            }
          }
        },

        # ── Row 4, col 0-5: Instance Count ─────────────────────────────────
        {
          xPos   = 0
          yPos   = 12
          width  = 6
          height = 4
          widget = {
            title = "Container Instance Count"
            xyChart = {
              dataSets = [{
                timeSeriesQuery = {
                  timeSeriesFilter = {
                    filter = join(" ", [
                      "metric.type=\"run.googleapis.com/container/instance_count\"",
                      "resource.type=\"cloud_run_revision\"",
                    ])
                    aggregation = {
                      alignmentPeriod    = "60s"
                      perSeriesAligner   = "ALIGN_MEAN"
                      crossSeriesReducer = "REDUCE_SUM"
                      groupByFields      = ["resource.labels.service_name"]
                    }
                  }
                }
                plotType       = "LINE"
                legendTemplate = "$${resource.labels.service_name}"
              }]
              yAxis = { label = "instances", scale = "LINEAR" }
            }
          }
        },

        # ── Row 4, col 6-11: Container Startup Latency ─────────────────────
        {
          xPos   = 6
          yPos   = 12
          width  = 6
          height = 4
          widget = {
            title = "Container Startup Latency — p50 / p95 (ms)"
            xyChart = {
              dataSets = [
                {
                  timeSeriesQuery = {
                    timeSeriesFilter = {
                      filter = join(" ", [
                        "metric.type=\"run.googleapis.com/container/startup_latencies\"",
                        "resource.type=\"cloud_run_revision\"",
                      ])
                      aggregation = {
                        alignmentPeriod    = "60s"
                        perSeriesAligner   = "ALIGN_DELTA"
                        crossSeriesReducer = "REDUCE_PERCENTILE_50"
                        groupByFields      = ["resource.labels.service_name"]
                      }
                    }
                  }
                  plotType       = "LINE"
                  legendTemplate = "$${resource.labels.service_name} p50"
                },
                {
                  timeSeriesQuery = {
                    timeSeriesFilter = {
                      filter = join(" ", [
                        "metric.type=\"run.googleapis.com/container/startup_latencies\"",
                        "resource.type=\"cloud_run_revision\"",
                      ])
                      aggregation = {
                        alignmentPeriod    = "60s"
                        perSeriesAligner   = "ALIGN_DELTA"
                        crossSeriesReducer = "REDUCE_PERCENTILE_95"
                        groupByFields      = ["resource.labels.service_name"]
                      }
                    }
                  }
                  plotType       = "LINE"
                  legendTemplate = "$${resource.labels.service_name} p95"
                },
              ]
              yAxis = { label = "ms", scale = "LINEAR" }
            }
          }
        },

      ]
    }
  })

  depends_on = [google_project_service.apis]
}
