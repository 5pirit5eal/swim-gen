#!/bin/bash
set -e

if [ -f .env ]; then
  source .env
fi

run() {
  go run ./main.go
}

validate() {
  go vet ./...
}

format() {
  go fmt ./...
}

test() {
  go test -v ./...
}

docs() {
  swag init
}

create-identity-token() {
  gcloud auth print-identity-token
}

authenticate() {
  gcloud auth login --update-adc --no-launch-browser
}

activate() {
  gcloud config configurations activate "$PROJECT_ID"
  authenticate
  gcloud auth application-default set-quota-project "$PROJECT_ID"
  echo "SUCCESS: GOOGLE CLOUD CONFIGURATION ACTIVATED"
}

setup-gcloud() {
  echo "--- SETTING UP LOCAL GOOGLE CLOUD SDK CONFIGURATION ---"
  gcloud config configurations create "$PROJECT_ID"
  activate
  gcloud config set project "$PROJECT_ID"
  gcloud config set compute/region "$REGION"
}

docker-run() {
  local container_id=$1
  local port=${2:-"8080"}
  # local background=${3:-"-d"}
  docker run \
    -v ~/.config/gcloud/application_default_credentials.json:/gcp/creds.json \
    -p $port:8080 -e PORT=8080 \
    -e GOOGLE_APPLICATION_CREDENTIALS=/gcp/creds.json \
    -e GOOGLE_CLOUD_PROJECT="$PROJECT_ID" \
    --env-file .env \
    $background \
    -i $container_id
}

docker-build-and-run() {
  docker build .
  docker-run $(docker images --format "{{.ID}}" | head -n 1)
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
