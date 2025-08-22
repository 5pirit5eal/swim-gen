locals {
  backend_env_variables = {
    _PROJECT_ID       = var.project_id
    _REGION           = var.region
    _MODEL            = var.model
    _SMALL_MODEL      = var.small_model
    _EMBEDDING_NAME   = var.embedding_name
    _EMBEDDING_MODEL  = var.embedding_model
    _EMBEDDING_SIZE   = var.embedding_size
    _DB_NAME          = var.dbname
    _DB_INSTANCE      = google_sql_database_instance.main.connection_name
    _DB_USER          = var.dbuser
    _DB_PASS_LOCATION = data.google_secret_manager_secret_version_access.dbpassword_user.id
    _PORT             = var.port
    _LOG_LEVEL        = var.log_level
    _BUCKET_NAME      = google_storage_bucket.exported_pdfs.name
    _SIGNING_SA       = google_service_account.pdf_export_sa.email
  }
  mcp_env_variables = {
    _PROJECT_ID       = var.project_id
    _REGION           = var.region
    _SERVICE_CPU      = var.service_cpu
    _SERVICE_MEMORY   = var.service_memory
    _SERVICE_TIMEOUT  = var.service_timeout
    _LOG_LEVEL        = var.log_level
    _SWIM_RAG_API_URL = var.backend_url
  }
  frontend_env_variables = {
    _PROJECT_ID = var.project_id
    _REGION     = var.region
    _BFF_URL    = var.bff_url
    _PORT       = var.port
  }
  bff_env_variables = {
    _PROJECT_ID   = var.project_id
    _REGION       = var.region
    _BACKEND_URL  = var.backend_url
    _FRONTEND_URL = var.domain_url
    _PORT         = var.port
    _LOG_LEVEL    = var.log_level
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

# github connection using the cloudbuildv2
resource "google_cloudbuildv2_connection" "github" {
  project  = var.project_id
  location = "europe-west1" // Needs to be europe-west1 for Cloud Build compatibility
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

# cloud build triggers for mcp server
resource "google_cloudbuild_trigger" "swim_rag_mcp_server_pr_main" {
  name               = "swim-rag-mcp-server-pr-main"
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

  included_files = ["mcp-server/**"]

  substitutions = local.mcp_env_variables
  tags          = ["mcp-server", "PR", "swim-rag", "main"]

  filename = "mcp-server/main-pr.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "swim_rag_mcp_server_release" {
  name               = "swim-rag-mcp-server-release"
  description        = "Trigger for swim-rag release from main branch"
  service_account    = google_service_account.cloud_build_sa.id
  location           = "europe-west1"
  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  repository_event_config {
    repository = google_cloudbuildv2_repository.swim_rag.id
    push {
      branch = "main"
    }
  }

  substitutions = merge(local.mcp_env_variables, {
    _AR_REPO_NAME    = google_artifact_registry_repository.docker.name
    _SERVICE_ACCOUNT = google_service_account.swim_gen_backend_sa.email
  })

  included_files = ["mcp-server/**"]

  tags = ["mcp-server", "PR", "swim-rag", "main"]

  filename = "mcp-server/release.cloudbuild.yaml"
}

# cloud build triggers for backend server
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

  included_files = ["backend/**"]

  substitutions = local.backend_env_variables
  tags          = ["backend", "PR", "swim-rag", "main"]

  filename = "backend/main-pr.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "swim_rag_backend_release" {
  name               = "swim-rag-backend-release"
  description        = "Trigger for swim-rag release from main branch"
  service_account    = google_service_account.cloud_build_sa.id
  location           = "europe-west1"
  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  repository_event_config {
    repository = google_cloudbuildv2_repository.swim_rag.id
    push {
      branch = "main"
    }
  }

  included_files = ["backend/**"]

  substitutions = merge(local.backend_env_variables, {
    _AR_REPO_NAME    = google_artifact_registry_repository.docker.name
    _SERVICE_ACCOUNT = google_service_account.swim_gen_backend_sa.email
  })

  tags = ["backend", "PR", "swim-rag", "main"]

  filename = "backend/release.cloudbuild.yaml"
}

# cloud build triggers for frontend server
resource "google_cloudbuild_trigger" "swim_rag_frontend_pr_main" {
  name               = "swim-rag-frontend-pr-main"
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

  included_files = ["frontend/**"]

  substitutions = local.frontend_env_variables
  tags          = ["frontend", "PR", "swim-rag", "main"]

  filename = "frontend/main-pr.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "swim_rag_frontend_release" {
  name               = "swim-rag-frontend-release"
  description        = "Trigger for swim-rag release from main branch"
  service_account    = google_service_account.cloud_build_sa.id
  location           = "europe-west1"
  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  repository_event_config {
    repository = google_cloudbuildv2_repository.swim_rag.id
    push {
      branch = "main"
    }
  }

  included_files = ["frontend/**"]

  substitutions = merge(local.frontend_env_variables, {
    _AR_REPO_NAME    = google_artifact_registry_repository.docker.name
    _SERVICE_ACCOUNT = google_service_account.swim_gen_frontend_sa.email
  })

  tags = ["frontend", "PR", "swim-rag", "main"]

  filename = "frontend/release.cloudbuild.yaml"
}

# cloud build triggers for bff server
resource "google_cloudbuild_trigger" "swim_rag_bff_pr_main" {
  name               = "swim-rag-bff-pr-main"
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

  included_files = ["bff/**"]

  substitutions = local.bff_env_variables
  tags          = ["bff", "PR", "swim-rag", "main"]

  filename = "bff/main-pr.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "swim_rag_bff_release" {
  name               = "swim-rag-bff-release"
  description        = "Trigger for swim-rag release from main branch"
  service_account    = google_service_account.cloud_build_sa.id
  location           = "europe-west1"
  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  repository_event_config {
    repository = google_cloudbuildv2_repository.swim_rag.id
    push {
      branch = "main"
    }
  }

  included_files = ["bff/**"]

  substitutions = merge(local.bff_env_variables, {
    _AR_REPO_NAME    = google_artifact_registry_repository.docker.name
    _SERVICE_ACCOUNT = google_service_account.swim_gen_frontend_sa.email
  })

  tags = ["bff", "PR", "swim-rag", "main"]

  filename = "bff/release.cloudbuild.yaml"
}
