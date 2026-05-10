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
  --build-arg VITE_SUPABASE_ANON_KEY="$VITE_SUPABASE_ANON_KEY" \
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

**PostGres Plugins**: Activate pgvector and pg_cron on your supase databases manually so that the system can use these dependencies.

#### Telemetry & Billing (One-Time)

After applying the `0-infra` stage with the telemetry and billing resources, complete these manual steps.

**Architecture overview:**

```
Cloud Run services (stdout/stderr + request logs)
  └─→ Cloud Logging log bucket (swim-gen-cloud-run, 1825-day retention)
        └─→ Linked BigQuery dataset (swim_gen_<env>_logs) — read-only, _AllLogs virtual view
              └─→ Analytics views in swim_gen_<env>_analytics dataset
                    └─→ Looker Studio dashboard

GCP Billing export
  └─→ BigQuery dataset (swim_gen_billing, prod only)
        └─→ v_costs view in swim_gen_prod_analytics
```

**Step 1 — Billing Export to BigQuery** (prod only):

- Go to GCP Console → Billing → Billing export → BigQuery export
- Enable **Detailed usage cost** export
- Select project: `swim-gen-prod`, dataset: `swim_gen_billing`
- GCP will begin populating a table named `gcp_billing_export_v1_<BILLING_ACCOUNT_ID>`
- Note: data appears with a ~1-day delay; the `v_costs` view will return no rows until the first export lands

**Step 2 — Trace Linked Dataset** (both environments, optional):

The `_Trace` linked dataset enables querying Cloud Trace spans directly from BigQuery. There is no Terraform resource for this — create it once with:

```bash
# Dev
gcloud beta observability buckets datasets links create swim_gen_dev_traces \
  --bucket=swim-gen-cloud-run \
  --location=europe-west4 \
  --project=rubenschulze-sandbox

# Prod
gcloud beta observability buckets datasets links create swim_gen_prod_traces \
  --bucket=swim-gen-cloud-run \
  --location=europe-west4 \
  --project=swim-gen-prod
```

**Step 3 — Looker Studio Dashboard:**

See the [Looker Studio Dashboard](#looker-studio-dashboard) section below.

---

#### Looker Studio Dashboard

Go to [lookerstudio.google.com](https://lookerstudio.google.com) → **Create → Report**.

##### Add data sources

Add one BigQuery data source per view. For each: **Add data → BigQuery → My projects → select project → `swim_gen_<env>_analytics`** → select the view.

| Data source name | View | Project |
|---|---|---|
| `swim-gen: Monthly Active Users` | `v_monthly_active_users` | dev or prod |
| `swim-gen: Total Users` | `v_total_users` | dev or prod |
| `swim-gen: Request Volume` | `v_request_volume` | dev or prod |
| `swim-gen: Request Latency` | `v_request_latency` | dev or prod |
| `swim-gen: Error Rate` | `v_error_rate` | dev or prod |
| `swim-gen: Costs` | `v_costs` | prod only |

> All views return no rows until traffic has flowed through the services and logs have been ingested. This is expected behaviour.

##### Date range control

Add a **Date range control** (Insert menu → Date range control). Set the default to **Last 12 months**. Connect it to all charts that have a `day` or `month` dimension — Looker Studio filters automatically when the dimension type is **Date** or **Year Month**.

For charts using `v_monthly_active_users`, set the date dimension to `month` with type **Year Month**.

##### Page 1 — Users

| Chart | Type | Data source | Dimension | Metric | Notes |
|---|---|---|---|---|---|
| Monthly Active Users | Time series | Monthly Active Users | `month` | `active_users` | Add `service` as breakdown dimension |
| All-time unique users | Scorecard | Total Users | — | `unique_users` | Add filter: `period = all_time` |
| Yearly unique users | Bar chart | Total Users | `year` | `unique_users` | Add filter: `period = yearly` |

##### Page 2 — Performance

| Chart | Type | Data source | Dimension | Metric | Notes |
|---|---|---|---|---|---|
| Request Volume (total) | Time series | Request Volume | `day` | `request_count` | Add filter `route = _total`; use `service` as breakdown |
| Request Volume by Route | Bar chart | Request Volume | `route` | `request_count` | Add filter `route != _total`; add `service` filter control |
| Latency p50/p95/p99 (total) | Time series | Request Latency | `day` | `p50_ms`, `p95_ms`, `p99_ms` | Add filter `route = _total`; add `service` filter control |
| Latency by Route | Table | Request Latency | `route`, `service` | `p50_ms`, `p95_ms`, `p99_ms` | Add filter `route != _total` |
| Error Rate (total) | Time series | Error Rate | `day` | `error_rate_pct` | Add filter `route = _total`; add `service` filter control |
| Error Rate by Route | Table | Error Rate | `route`, `service` | `error_rate_pct`, `server_error_rate_pct`, `total_requests` | Add filter `route != _total` |

> The `route` field uses the value `_total` as a sentinel for the rolled-up daily total across all routes. Filter on `route = _total` for aggregate charts and `route != _total` for per-route breakdowns.

##### Page 3 — Costs (prod only)

Connect these charts to the prod project's `swim_gen_prod_analytics` dataset.

| Chart | Type | Data source | Dimension | Metric | Notes |
|---|---|---|---|---|---|
| Daily Cost | Stacked area | Costs | `day` | `daily_cost` | Use `gcp_service` as breakdown dimension |
| Cost by SKU | Bar chart | Costs | `sku` | `daily_cost` | Sort descending by `daily_cost` |
| Total cost (period) | Scorecard | Costs | — | `SUM(daily_cost)` | Respects the date range control |

##### Filter controls (report-level)

Add these to the report header so they apply across all pages:

- **Service** — dropdown filter on the `service` field
- **Date range** — already added above

Add this to Page 2 only:

- **Route** — dropdown filter on the `route` field; set a default exclusion filter of `route != _total` on all per-route charts to avoid the rollup row appearing in breakdowns

##### Switching to prod

Once validated in dev, clone the report (**File → Make a copy**), then for each data source swap the project from `rubenschulze-sandbox` / `swim_gen_dev_analytics` to `swim-gen-prod` / `swim_gen_prod_analytics`, and add Page 3 (costs).

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
