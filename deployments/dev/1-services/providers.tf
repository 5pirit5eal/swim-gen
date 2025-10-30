terraform {
  backend "gcs" {
    bucket = "swim-gen-state-dev"
    prefix = "tofu/swim-gen-services"
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

provider "github" {
  token = var.github_token
  owner = var.github_owner
}
