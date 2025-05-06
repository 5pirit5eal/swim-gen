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

# github connection using the cloudbuildv2
resource "google_cloudbuildv2_connection" "github" {
  project  = var.project_id
  location = "europe-west1"
  name     = "swim-rag-github-connection"

  github_config {
    app_installation_id = var.github_app_installation_id
    authorizer_credential {
      oauth_token_secret_version = google_secret_manager_secret_version.github_token_secret_version.id
    }
  }
  depends_on = [
    google_secret_manager_secret_iam_member.cloud_build_sa_secret_access,
    google_secret_manager_secret_version.github_token_secret_version,
    google_secret_manager_secret_iam_member.cloud_build_sa_github_token_access
  ]
}

resource "google_cloudbuildv2_repository" "swim_rag" {
  project           = var.project_id
  location          = "europe-west1"
  name              = "swim-rag"
  parent_connection = google_cloudbuildv2_connection.github.name
  remote_uri        = "${var.github_uri}.git"
}

# cloud build triggers
resource "google_cloudbuild_trigger" "swim_rag_backend_pr_main" {
  name               = "swim-rag-backend-pr-main"
  description        = "Trigger for swim-rag PR to main branch"
  service_account    = google_service_account.cloud_build_sa.id
  location           = "europe-west1"
  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  repository_event_config {
    repository = google_cloudbuildv2_repository.swim_rag.id
    pull_request {
      branch = "main"
    }
  }

  filename = "backend/main-pr.cloudbuild.yaml"
}