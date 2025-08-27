terraform {
  backend "gcs" {
    bucket = "rubenschulze-sandbox-state"
    prefix = "tofu/swim-rag"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

data "google_project" "project" {
  project_id = var.project_id
}