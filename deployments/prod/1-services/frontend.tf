locals {
  frontend_env_variables = {
    PROJECT_ID       = var.project_id
    REGION           = var.region
    VITE_APP_API_URL = google_cloud_run_v2_service.bff.uri
  }
  records = {
    for type, records in {
      for r in google_cloud_run_domain_mapping.frontend_domain_mapping.status[0].resource_records : r.type => r.rrdata...
    } : type => { rrdatas = records }
  }
}

data "google_artifact_registry_docker_image" "frontend_image" {
  location      = data.google_artifact_registry_repository.docker.location
  repository_id = data.google_artifact_registry_repository.docker.repository_id
  image_name    = "swim-gen-frontend:${var.version_tag}"
}

resource "google_cloud_run_v2_service" "frontend" {
  name     = "swim-gen-frontend"
  location = var.region

  # gcloud command used --no-allow-unauthenticated (public ingress but no allUsers binding)
  # So expose all ingress, rely on IAM to restrict.
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true

  # Only allow authenticated invocations from other services with the "run.invoker" role
  # (e.g. the frontend service).
  # This is enforced via IAM bindings in the frontend.tf file.
  # See: https://cloud.google.com/run/docs/securing/service-identity#granting_other_identities_access_to_your_service
  deletion_protection = false

  custom_audiences = [var.domain_url]

  # Set the number of maximum instances to control costs
  scaling {
    max_instance_count = 5
  }

  template {
    service_account                  = var.iam.swim_gen_frontend.email
    max_instance_request_concurrency = 200
    timeout                          = "60s"
    containers {
      image = data.google_artifact_registry_docker_image.frontend_image.self_link
      resources {
        limits = {
          cpu    = 1
          memory = "512Mi"
        }
        cpu_idle          = true
        startup_cpu_boost = true
      }
      liveness_probe {
        http_get {
          path = "/health"
          port = 8080
        }
        initial_delay_seconds = 5
        period_seconds        = 10
        timeout_seconds       = 2
        failure_threshold     = 3
      }

      dynamic "env" {
        for_each = local.frontend_env_variables
        content {
          name  = env.key
          value = tostring(env.value)
        }
      }
    }
  }

  client     = "terraform"
  depends_on = [google_cloud_run_v2_service.backend, google_cloud_run_v2_service.frontend]
}

# Cloud domain mapping
resource "google_cloud_run_domain_mapping" "frontend_domain_mapping" {
  location = google_cloud_run_v2_service.frontend.location
  name     = var.domain_url

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.frontend.name
  }
}

# DNS records for the domain mapping
resource "google_dns_record_set" "frontend_dns_records" {
  for_each     = local.records
  name         = data.google_dns_managed_zone.swim_gen_zone.dns_name
  managed_zone = data.google_dns_managed_zone.swim_gen_zone.name
  type         = each.key
  ttl          = 300

  rrdatas    = each.value.rrdatas
  depends_on = [google_cloud_run_domain_mapping.frontend_domain_mapping]
}

data "google_dns_managed_zone" "swim_gen_zone" {
  name = "swim-gen-com"
}

# Site URL Github Actions Env Variable
data "github_repository" "swim_gen_repo" {
  full_name = "${var.github_owner}/${var.github_repository}"
}

resource "github_actions_environment_variable" "prod_site_url" {
  repository    = data.github_repository.swim_gen_repo.name
  environment   = "prod"
  variable_name = "VITE_SITE_URL"
  value         = "https://${var.domain_url}"
}
