# This file contains the variables for the Terraform configuration.
project_id = "rubenschulze-sandbox"
region     = "europe-west4"
apis = [
  "cloudresourcemanager.googleapis.com", # Terraform Backend
  "run.googleapis.com",                  # Cloud Run
  "container.googleapis.com",            # Container Registry
  "artifactregistry.googleapis.com",     # Artifact Registry
  "aiplatform.googleapis.com",           # Vertex AI
  "monitoring.googleapis.com",           # Cloud Monitoring
  "cloudtrace.googleapis.com",           # Cloud Trace
  "bigquery.googleapis.com",             # BigQuery
  "logging.googleapis.com",              # Cloud Logging
  "telemetry.googleapis.com",            # Telemetry
  "billingbudgets.googleapis.com"        # Billing Budgets
]
dbusers = {
  backend  = "coach"
  frontend = "swimmer"
}
outputs_location = "../0-config"
supabase = {
  organization_id = "rbrfvltmypsayebvplbb"
  name            = "swim-gen-dev"
  region          = "eu-central-1"
}
