# This file contains the variables for the Terraform configuration.
project_id = "rubenschulze-sandbox"
region     = "europe-west4"
apis = [
  "cloudresourcemanager.googleapis.com", # Terraform Backend
  "run.googleapis.com",                  # Cloud Run
  "container.googleapis.com",            # Container Registry
  "artifactregistry.googleapis.com",     # Artifact Registry
  "aiplatform.googleapis.com"            # Vertex AI
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
