# This file contains the variables for the Terraform configuration.
project_id = "swim-gen-prod"
region     = "europe-west4"
apis = [
  "cloudresourcemanager.googleapis.com",
  "run.googleapis.com",
  "container.googleapis.com",
  "compute.googleapis.com",
  "artifactregistry.googleapis.com",
  "aiplatform.googleapis.com",
  "monitoring.googleapis.com",
  "cloudtrace.googleapis.com",
  "bigquery.googleapis.com",
  "logging.googleapis.com",
  "telemetry.googleapis.com",
  "billingbudgets.googleapis.com"
]
dbusers = {
  backend  = "coach"
  frontend = "swimmer"
}

supabase = {
  organization_id = "rbrfvltmypsayebvplbb"
  name            = "swim-gen-prod"
  region          = "eu-central-1"
}
outputs_location = "../0-config"
