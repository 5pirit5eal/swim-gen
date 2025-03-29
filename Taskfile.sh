#!/bin/bash
set -e

source .env

run() {
  echo "Running the application..."
  go run ./cmd/swim-rag/swim-rag.go
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
    -e IMAGE_BUCKET="$IMAGE_BUCKET" \
    -e PROJECT_ID="$PROJECT_ID" \
    -e LOG_LEVEL=DEBUG \
    $background \
    -i $container_id
}

start-proxy() {
    mkdir -p /tmp/cloudsql
    chmod 777 /tmp/cloudsql
    sudo docker run -d -v /tmp/cloudsql:/cloudsql \
      -v ~/.config/gcloud/application_default_credentials.json:/gcp/creds.json \
      gcr.io/cloud-sql-connectors/cloud-sql-proxy:2.15.2 --unix-socket=/cloudsql \
      --credentials-file /gcp/creds.json "$DB_INSTANCE"
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