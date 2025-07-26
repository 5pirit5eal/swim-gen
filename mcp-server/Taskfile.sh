#!/bin/bash
set -e

source .config.env
# //////////////////////////////////////////////////////////////////////////////
# START tasks

run() {
  uv run fastmcp run src/swim_rag_mcp/main.py \
    --transport http 
}

dev() {
  uv run fastmcp dev --server-spec src/swim_rag_mcp/main.py
}

docker-run() {
  local container_id=$1
  local port=${2:-"8080"}
  # local background=${3:-"-d"}
  docker run \
    -v ~/.config/gcloud/application_default_credentials.json:/gcp/creds.json \
    -p $port:8080 \
    -e GOOGLE_APPLICATION_CREDENTIALS=/gcp/creds.json \
    -e GOOGLE_CLOUD_PROJECT="$PROJECT_ID" \
    --env-file .config.env \
    $background \
    -i $container_id
}

install() {
  uv sync --group dev
}

format() {
  local files=${1:-'.'}
  uv run ruff format "$files"
  uv run ruff check --select I --fix "$files"
}

validate() {
  local files=${1:-'.'}
  local mypy_cache=${2:-'.mypy_cache'}
  uv run ruff check "$files"
  uv run ruff format "$files" --diff
  uv run ruff check --select I "$files"
  mkdir -p "$mypy_cache"
  echo "Running mypy..."
  uv run mypy "$files" --cache-dir "$mypy_cache"
  
  echo "Security check with bandit..."
  uv run bandit -c pyproject.toml -ll -x ./.venv/ -r "$files"
}

integration-test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  local port=${2:-8080}
  # Start emulator
  gcloud emulators firestore start --host-port="localhost:$port" > /dev/null 2>&1 &
  export FIRESTORE_EMULATOR_HOST="localhost:$port"
  sleep 10 # Wait for emulator to start
  uv run --group dev pytest tests/integration_tests/test_* "$*"
  result=$?
  # Stop emulator
  curl -d '' "localhost:$port/shutdown" > /dev/null 2>&1 &
  sleep 1
  return $result
}

unit-test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  uv run --group dev pytest tests/unit_tests/test_* "$*"
  result=$?
  sleep 1
  return $result
}

system-test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  uv run --group dev pytest tests/system_tests/test_* "$*"
  result=$?
  sleep 1
  return $result
}

create-identity-token() {
  gcloud auth print-identity-token
}

authenticate() {
  gcloud auth login --update-adc --no-launch-browser
}

activate() {
  # ADAPT IF WANTING TO USE E.G. CONDA
  gcloud config configurations activate "$PROJECT_ID"
  authenticate
  gcloud auth application-default set-quota-project "$PROJECT_ID"
  echo "SUCCESS: GOOGLE CLOUD CONFIGURATION ACTIVATED"
}

setup-gcloud() {
  local setup_wif=${1:-"false"}
  local google_account=${2:-""}
  echo "--- SETTING UP LOCAL GOOGLE CLOUD SDK CONFIGURATION ---"
  gcloud config configurations create "$PROJECT_ID"
  if [ -n "$google_account" ]; then
    gcloud config set account "$google_account"
  fi
  activate
  gcloud config set project "$PROJECT_ID"
  gcloud config set compute/region "$REGION"
  gcloud components install cloud-firestore-emulator

  if [ "$setup_wif" = "true" ]; then
    setup-wif
  fi
}

setup-wif() {
  # SOURCE: https://github.com/google-github-actions/auth#indirect-wif
  # info: all other env variables are imported via "source config.env" (line 5)
  echo "--- SETTING UP GOOGLE CLOUD INFRASTRUCTURE ---"
  echo "creating service account: $CLOUD_RUN_SERVICE_ACCOUNT"
  gcloud iam service-accounts create "$CLOUD_RUN_SERVICE_ACCOUNT" \
    --project "$PROJECT_ID"

  echo "writing 'SERVICE_ACCOUNT=$CLOUD_RUN_SERVICE_ACCOUNT' to config.env"
  echo "SERVICE_ACCOUNT=\"$CLOUD_RUN_SERVICE_ACCOUNT\"" >> config.env

  echo "creating workload-identity-pool in project $PROJECT_ID"
  gcloud iam workload-identity-pools create "github" \
  --project="$PROJECT_ID" \
  --location="global" \
  --display-name="GitHub Actions Pool"

  echo "checking for created workload-identity-pools in project $PROJECT_ID"
  WIF_POOL_NAME=$(gcloud iam workload-identity-pools describe "github" \
  --project="$PROJECT_ID" \
  --location="global" \
  --format="value(name)")
  echo "success: workload-identity-pool $WIF_POOL_NAME exists"

  REPO_NAME=$(awk -F' *= *' 'NR==2 && /name *=/ {gsub(/"/, "", $2); print $2}' pyproject.toml)

  echo "adding iam-policy-binding for $CLOUD_RUN_SERVICE_ACCOUNT in project $PROJECT_ID"
  uv run "gcloud iam service-accounts add-iam-policy-binding $CLOUD_RUN_SERVICE_ACCOUNT" \
    --project="$PROJECT_ID" \
    --role="roles/iam.workloadIdentityUser" \
    --member="principalSet://iam.googleapis.com/$WIF_POOL_NAME/attribute.repository/$GITHUB_ORG/$REPO_NAME"

  WIF_PROVIDER_NAME=$(uv run gcloud iam workload-identity-pools providers describe "$REPO_NAME" \
  --project="$PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="github" \
  --format="value(name)")

  echo "checking for workload-identity-provider name in project $PROJECT_ID"
  echo "success: workload-identity-provider $WIF_PROVIDER_NAME was retrieved"
  echo "writing 'WIF_PROVIDER_NAME=$WIF_PROVIDER_NAME' to config.env"
  echo "WIF_PROVIDER_NAME=\"$WIF_PROVIDER_NAME\"" >> config.env
  echo "SUCCESS: SETTING UP GOOGLE CLOUD INFRASTRUCTURE IS COMPLETE"
}

# END tasks
# //////////////////////////////////////////////////////////////////////////////
# START CI/CD tasks


install-automation() { # install using uv
  npm install -g semantic-release @semantic-release/exec @adesso-gcc/semantic-release-config
  uv sync --frozen
}

prepare-release() {
  # Bump the version in pyproject.toml
  echo "version = $1"

  # Get version number from version tag
  local version
  version=`echo "$1" | cut -d'v' -f2`
  echo "py = $version"

  # Use perl to replace the version in pyproject.toml
  perl -pi -e 's/^version = .*/version = "'$version'"/' pyproject.toml

  # Update lock file with the new version
  uv lock
}

release() {
  npx semantic-release $*
}

test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  set +e
  unit-test "$*"
  unit_test_status=$?
  integration-test "$*"
  integration_test_status=$?
  # system-test "$*"
  set -e
  if [ $unit_test_status -ne 0 ] || [ $integration_test_status -ne 0 ]; then
    if [ $unit_test_status -ne 0 ]; then
      echo "Unit tests failed"
    fi
    if [ $integration_test_status -ne 0 ]; then
      echo "Integration tests failed"
    fi
    return 1
  fi
}

deploy() {
  local project_id=${1:-$PROJECT_ID}
  local region=${2:-$REGION}
  local pw=${3:-$PW_SECRET_NAME}
  local sa=${4:-$CLOUD_RUN_SERVICE_ACCOUNT}
  local platform=${5:-linux/amd64}

  echo "... fetching project variables ..."
  # Add other variables if necessary
  NAME=$(awk -F' *= *' 'NR==2 && /name *=/ {gsub(/"/, "", $2); print $2}' pyproject.toml)
  VERSION=$(awk -F' *= *' '/version *=/ {gsub(/"/, "", $2); print $2}' pyproject.toml)
  DATETIME=$(date +"%y-%m-%d-%H%M%S")

  IMAGE_NAME="$region-docker.pkg.dev/$project_id/docker/$NAME"
  IMAGE_TAG="$region-docker.pkg.dev/$project_id/docker/$NAME:$VERSION-$DATETIME"

  echo "... building docker image ..."
  gcloud auth configure-docker "$region-docker.pkg.dev"
  docker build --platform "$platform" --tag "$IMAGE_TAG" .

  echo "... pushing image to artifact registry ..."
  docker push "$IMAGE_TAG"
  gcloud artifacts docker tags add "$IMAGE_TAG" "$IMAGE_NAME":latest
  gcloud artifacts docker tags add "$IMAGE_TAG" "$IMAGE_NAME:$VERSION"

  gcloud run deploy "$NAME" \
  --image "$IMAGE_TAG" \
  --platform managed \
  --project="$project_id" \
  --region="$region" \
  --cpu "$SERVICE_CPU" \
  --memory "$SERVICE_MEMORY" \
  --timeout "$SERVICE_TIMEOUT" \
  --service-account "$sa" \
  --set-env-vars=PROJECT_ID="$project_id",REGION="$region",PW_SECRET_NAME="$pw",LOG_LEVEL="$LOG_LEVEL"
}

# END CI/CD tasks
# //////////////////////////////////////////////////////////////////////////////
help() {
  echo "Usage: ./Taskfile.sh [task]"
  echo
  echo "Available tasks:"
  echo "  run                           Run the application locally."
  echo "  docker-run                    Run the application in a previously built Docker container."
  echo "  format                        Format the code using ruff."
  echo "  validate                      Perform code linting and formatting using ruff and mypy."
  echo "  install                       Install development dependencies."
  echo "  integration-test              Run integration tests using pytest."
  echo "  unit-test                     Run unit tests using pytest."
  echo "  system-test                   Run system tests using pytest."
  echo "  create-identity-token         Create an identity token for external request authentication."
  echo "  authenticate                  Authenticate to Google Cloud."
  echo "  activate                      Activate Google Cloud configuration."
  echo "  setup-gcloud                  Set up the Google Cloud settings (and Workload Identity Federation)."
  echo "  setup-wif                     Set up Workload Identity Federation."
  echo
  echo "Automation tasks:"
  echo "  install-automation            Install dependencies using uv."
  echo "  prepare-release               Prepare the release by bumping the version."
  echo "  release                       Perform the release using semantic-release."
  echo "  test                          Run unit tests."
  echo "  deploy                        Deploy application to Google Cloud Run."
  echo
  echo "If no task is provided, the default is to run the application."
}

# Check if the provided argument matches any of the functions
if [ -n "$1" ] && ! declare -f "$1" > /dev/null; then
  echo "Error: Unknown task '$1'"
  echo
  help  # Show help if the task is not recognized
  exit 1
fi

# Run application if no argument is provided
"${@:-run}"