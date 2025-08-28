locals {
  github_env_variables = {
    PROJECT_ID   = var.project_id
    REGION       = var.region
    AR_REPO_NAME = google_artifact_registry_repository.docker.repository_id
    WIF_PROVIDER = "projects/${data.google_project.project.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.github.workload_identity_pool_id}/providers/${google_iam_workload_identity_pool_provider.github.workload_identity_pool_provider_id}"
    WIF_SA       = google_service_account.github_actions_sa.email
  }
}


resource "google_artifact_registry_repository" "docker" {
  location               = var.region
  repository_id          = "docker"
  description            = "Docker repository for my sandbox projects using cloud run."
  format                 = "DOCKER"
  cleanup_policy_dry_run = false
  cleanup_policies {
    id     = "delete-untagged"
    action = "DELETE"
    condition {
      tag_state = "UNTAGGED"
    }
  }
  cleanup_policies {
    id     = "keep-new-untagged"
    action = "KEEP"
    condition {
      tag_state  = "UNTAGGED"
      newer_than = "7d"
    }
  }
  cleanup_policies {
    id     = "keep-tagged-release"
    action = "KEEP"
    condition {
      tag_state    = "TAGGED"
      tag_prefixes = ["release"]
    }
  }
}

# Workload Identity Federation
resource "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = "github"
  display_name              = "GitHub Actions"
  description               = "WIF pool for GitHub Actions OIDC"
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = "github"
  display_name                       = "GitHub OIDC"
  description                        = "Provider for token.actions.githubusercontent.com"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
    "attribute.ref"        = "assertion.ref"
    "attribute.actor"      = "assertion.actor"
    "attribute.workflow"   = "assertion.workflow"
  }

  # Limit to your GitHub org/user (optional but recommended)
  attribute_condition = "assertion.repository_owner == '${var.github_owner}'"
}

resource "google_service_account_iam_binding" "gh_actions_wif_repo" {
  service_account_id = google_service_account.github_actions_sa.name
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_owner}/${var.github_repository}"
  ]
}


# Github Actions Environments and Variables
data "github_repository" "swim_gen_repo" {
  full_name = "${var.github_owner}/${var.github_repository}"
}

data "github_repository_environments" "dev_environment" {
  repository = data.github_repository.swim_gen_repo.name
}

resource "github_repository_environment" "prod" {
  repository  = data.github_repository.swim_gen_repo.name
  environment = "prod"
}

resource "github_actions_environment_variable" "prod_project_id" {
  for_each      = local.github_env_variables
  repository    = data.github_repository.swim_gen_repo.name
  environment   = github_repository_environment.prod.environment
  variable_name = each.key
  value         = each.value
}

