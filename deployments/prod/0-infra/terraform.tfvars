# This file contains the variables for the Terraform configuration.
project_id = "swim-gen-prod"
region     = "europe-west4"
apis = [
  "cloudresourcemanager.googleapis.com",
  "run.googleapis.com",
  "sqladmin.googleapis.com",
  "container.googleapis.com",
  "compute.googleapis.com",
  "artifactregistry.googleapis.com",
  "aiplatform.googleapis.com"
]
dbname           = "swim-gen-db"
dbuser           = "swimmer"
dbtier           = "db-f1-micro"
outputs_location = "../0-config"
