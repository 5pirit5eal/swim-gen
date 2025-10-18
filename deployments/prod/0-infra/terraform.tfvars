# This file contains the variables for the Terraform configuration.
project_id = "swim-gen-prod"
region     = "europe-west4"
apis = [
  "cloudresourcemanager.googleapis.com",
  "run.googleapis.com",
  "container.googleapis.com",
  "compute.googleapis.com",
  "artifactregistry.googleapis.com",
  "aiplatform.googleapis.com"
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
