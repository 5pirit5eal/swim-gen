terraform {
  required_providers {
    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.0"
    }
  }
  backend "gcs" {
    bucket = "swim-gen-state-prod"
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

provider "supabase" {
  access_token = var.supabase_access_token
}

data "google_project" "project" {
  project_id = var.project_id
}
