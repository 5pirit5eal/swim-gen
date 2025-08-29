# This file contains the variables for the Terraform configuration.
project_id = "rubenschulze-sandbox"
region     = "europe-west4"
apis = [
  "cloudresourcemanager.googleapis.com",
  "run.googleapis.com",             # Cloud Run
  "sqladmin.googleapis.com",        # Cloud SQL
  "container.googleapis.com",       # Container Registry
  "compute.googleapis.com",         # Compute Engine (required for Cloud SQL)
  "artifactregistry.googleapis.com" # Artifact Registry
]
dbname           = "swim-gen-db"
dbuser           = "swimmer"
dbtier           = "db-f1-micro"
outputs_location = "../0-config"
