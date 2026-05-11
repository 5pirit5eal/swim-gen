# --------------------------------------------------------------------------
# BigQuery Dataset for Billing Export
# --------------------------------------------------------------------------

resource "google_bigquery_dataset" "billing" {
  dataset_id  = "swim_gen_billing"
  description = "GCP billing export data"
  project     = var.project_id
  location    = var.region

  # Keep data indefinitely
  # default_table_expiration_ms not set

  delete_contents_on_destroy = false

  depends_on = [google_project_service.apis]
}

# --------------------------------------------------------------------------
# Billing Export to BigQuery
# --------------------------------------------------------------------------
# Terraform does NOT support configuring billing data export to BigQuery.
# This must be done manually (one-time setup per billing account):
#
# Option A — GCP Console:
#   1. Go to Billing → Billing export → BigQuery export
#   2. Enable "Detailed usage cost" export
#   3. Select project: swim-gen-prod
#   4. Select dataset: swim_gen_billing
#
# Option B — gcloud CLI:
#   gcloud billing accounts describe <BILLING_ACCOUNT_ID>
#   (then configure via Console as CLI support is limited)
#
# After setup, GCP will create a table named
# `gcp_billing_export_v1_<BILLING_ACCOUNT_ID>` in the dataset.
# Use this table name in the BigQuery views (Part 8, v_costs).
# --------------------------------------------------------------------------
