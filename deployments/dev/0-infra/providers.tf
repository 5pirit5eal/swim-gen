terraform {
  required_version = ">= 1.0.4"
  required_providers {
    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.0"
    }
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = ">= 1.15.0"
    }
    github = {
      source = "integrations/github"
    }
  }
  backend "gcs" {
    bucket = "swim-gen-state-dev"
    prefix = "tofu/swim-gen-infra"
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
