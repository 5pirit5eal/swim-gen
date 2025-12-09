# Deployments

This directory contains the infrastructure-as-code (IaC) definitions for deploying the Swim RAG application to Google Cloud Platform (GCP).

If you want to know how to run the application, checkout [frontend/README.md](../frontend/README.md)

## Tooling

The infrastructure is managed using **OpenTofu**. All configurations are written in the Terraform language.

## Environments

There are two environments defined:

- `dev`: The development environment.
- `prod`: The production environment.

Each environment has two main stages:

- `0-infra`: Core infrastructure components like networking, IAM, and secrets.
- `1-services`: The application services (frontend, backend, BFF).

## Frontend build environment variables

The frontend build expects several environment variables to be available when building the Docker image or running the build locally. These values are used to populate the site's Impressum/contact information and other runtime-config.

Required variables

- VITE_IMPRESSUM_NAME
- VITE_IMPRESSUM_ADDRESS
- VITE_IMPRESSUM_CITY
- VITE_IMPRESSUM_PHONE
- VITE_IMPRESSUM_EMAIL
- VITE_SUPABASE_URL
- VITE_SUPABASE_ANON_KEY

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
  --build-arg VITE_SITE_URL="$VITE_SITE_URL" \
  --build-arg VITE_SUPABASE_URL="$VITE_SUPABASE_URL" \
  --build-arg VITE_SUPABASE_ANON_KEY="$VITE_SUPABASE_ANON_KEY"
  .
```

Notes

- If any of the variables are missing the build will receive empty values — consider adding a validation step in CI to fail early if a required secret is not set.
- For security, prefer GitHub Secrets (they are masked in logs) and avoid embedding sensitive information into final image layers.

## Deployment Process

### CI/CD with GitHub Actions

The project follows a trunk-based development model where all changes are committed directly to the `main` branch.

- **Validation**: On every commit to `main`, automated tests and linters run to validate the changes.
- **Deployment**: Deployments to the `dev` environment are triggered by adding a `.deploy dev` comment to a commit on the `main` branch. This action initiates a workflow that builds the services and deploys them to Cloud Run. Deployments to `prod` need to be triggered in the Github actions UI.

It is strongly recommended to rely on the CI/CD pipeline for all deployments to ensure consistency and safety.

### Manual Deployment

Manual deployments are discouraged but may be necessary for initial setup or specific maintenance tasks.

#### Prerequisites

Before you can run any OpenTofu commands, you must manually configure the following:

1. **Google Cloud Projects**: Create two google cloud projects, one for dev, one for prod. These projects are not created as part of the terraform configuration. You may adapt the code to create a project, but need to find another way of storing and providing the pre-required secrets.
2. **Domain & DNS**: Procure a domain name and configure its DNS settings to point to your GCP project. The specific DNS records will be output by the `0-infra` stage.
3. **Secrets**: Create and configure the necessary secrets in Google Secret Manager. This includes passwords for the Cloud SQL database users. These secrets must be created before the `0-infra` stage can be successfully applied, as it sets permissions on them. Check the data blocks in the dev and prod configuration for details.
4. **Google Auth**: Setup Google Identity Authentication following the [supabase tutorial](https://supabase.com/docs/guides/auth/social-login/auth-google). Set the site-url in the supabase project to the frontend url provided by cloud run or the domain of the app.
5. **Supabase Email Templates**: Configure the email templates in your Supabase project (Authentication -> Email Templates).
    - **Invite User**: Copy content from `deployments/supabase/templates/invite_user.html`
    - **Reset Password**: Copy content from `deployments/supabase/templates/reset_password.html`
    - **Confirm Signup**: Copy content from `deployments/supabase/templates/confirm_signup.html`

[!IMPORTANT]
The Sign-in with Google integration cannot be configured programmatically via Terraform/OpenTofu. You must set it up manually in the Google Cloud Console.

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

## Database Deployment

The projects relies on a postgres database for storing and retrieving training plans and user data. The database instance is provided in the `0-infra` stage.
