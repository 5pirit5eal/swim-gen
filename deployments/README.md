# Deployments

This directory contains the infrastructure-as-code (IaC) definitions for deploying the Swim RAG application to Google Cloud Platform (GCP).

## Tooling

The infrastructure is managed using **OpenTofu**. All configurations are written in the Terraform language.

## Environments

There are two environments defined:

- `dev`: The development environment.
- `prod`: The production environment.

Each environment has two main stages:

- `0-infra`: Core infrastructure components like networking, IAM, and secrets.
- `1-services`: The application services (frontend, backend, BFF, MCP server).

## Frontend build environment variables

The frontend build expects several environment variables to be available when building the Docker image or running the build locally. These values are used to populate the site's Impressum/contact information and other runtime-config.

Required variables

- VITE_IMPRESSUM_NAME
- VITE_IMPRESSUM_ADDRESS
- VITE_IMPRESSUM_CITY
- VITE_IMPRESSUM_PHONE
- VITE_IMPRESSUM_EMAIL

Where to set them

- CI (recommended): Add these as GitHub repository or environment secrets (exact names above). The GitHub Actions workflow `/.github/workflows/frontend-build.yaml` reads them from the `secrets` context and passes them as build-args to Docker.
- Locally: export them in your shell or add them to a local env file that your local build picks up (for example `frontend/.env.development` which is referenced by the workflow setup action).

Quick local example (zsh):

```bash
export VITE_IMPRESSUM_NAME="Acme GmbH"
export VITE_IMPRESSUM_ADDRESS="Street 1"
export VITE_IMPRESSUM_CITY="City"
export VITE_IMPRESSUM_PHONE="+49 000 000000"
export VITE_IMPRESSUM_EMAIL="hello@example.com"

# then build the frontend docker image from the repo root
cd frontend
docker build --build-arg VITE_IMPRESSUM_NAME="$VITE_IMPRESSUM_NAME" \
  --build-arg VITE_IMPRESSUM_ADDRESS="$VITE_IMPRESSUM_ADDRESS" \
  --build-arg VITE_IMPRESSUM_CITY="$VITE_IMPRESSUM_CITY" \
  --build-arg VITE_IMPRESSUM_PHONE="$VITE_IMPRESSUM_PHONE" \
  --build-arg VITE_IMPRESSUM_EMAIL="$VITE_IMPRESSUM_EMAIL" \
  .
```

Notes

- If any of the variables are missing the build will receive empty values — consider adding a validation step in CI to fail early if a required secret is not set.
- For security, prefer GitHub Secrets (they are masked in logs) and avoid embedding sensitive information into final image layers.

## Deployment Process

### CI/CD with GitHub Actions (Preferred Method)

All deployments are automated using GitHub Actions. Pushes and pull requests to the `main` and `dev` branches will trigger the corresponding workflows to plan and apply infrastructure changes.

The primary workflow for deployment is defined in `.github/workflows/tf-plan-apply.yaml`. This workflow is responsible for running `tofu plan` and `tofu apply` for each stage. Other workflows like `merge-dev.yaml` and `deploy-main.yaml` orchestrate the deployment process for the respective environments.

It is strongly recommended to rely on the CI/CD pipeline for all deployments to ensure consistency and safety.

### Manual Deployment

Manual deployments are discouraged but may be necessary for initial setup or specific maintenance tasks.

#### Prerequisites

Before you can run any OpenTofu commands, you must manually configure the following:

1. **Domain & DNS**: Procure a domain name and configure its DNS settings to point to your GCP project. The specific DNS records will be output by the `0-infra` stage.
2. **Secrets**: Create and configure the necessary secrets in Google Secret Manager. This includes passwords for the Cloud SQL database users. These secrets must be created before the `0-infra` stage can be successfully applied, as it sets permissions on them. Check the data blocks in the dev and prod configuration for details.

Check the data blocks for the relevant information that needs to be prepared before applying the configuration for the first time.

#### Initial Setup (One-Time)

The `0-infra` stage for each environment must be applied manually from your local machine **at least once**. This is because it provisions the core infrastructure that the CI/CD pipeline itself depends on, such as the service accounts and permissions used by GitHub Actions.

To run the initial setup for the `dev` environment:

```bash
# Ensure you are authenticated with GCP
gcloud auth application-default login

# Navigate to the infra directory
cd deployments/dev/0-infra

# Initialize OpenTofu
tofu init

# Plan the changes
tofu plan

# Apply the changes
tofu apply
```

Repeat this process for the `prod` environment if needed.

The prod environment is also designed to be used together with a custom domain. You should therefore ensure that the DNS zone is already setup.
Once the DNS records have been set by OpenTofu, they take around **15-30 mins** to take effect. The record setting relies on an already setup domain mapping, which requires you to use a targeted apply on the first run and then run it a second time.

```bash
tofu apply -target google_cloud_run_domain_mapping.frontend_domain_mapping
```

```bash
tofu apply
```

#### Subsequent Manual Runs

After the initial setup, you can run subsequent deployments manually if required.

**Example for `dev` environment `1-services` stage:**

```bash
# Ensure you are authenticated with GCP
gcloud auth application-default login

# Navigate to the services directory
cd deployments/dev/1-services

# Initialize OpenTofu
tofu init

# Plan the changes
tofu plan

# Apply the changes
tofu apply
```
