#!/bin/bash
set -euo pipefail

if [ -f .env ]; then
  # shellcheck source=/dev/null
  source .env
fi

run() {
  go run ./main.go
}

run-local() {
  go run ./main.go -env .env.local
}

validate() {
  go vet ./...
  golangci-lint run --path-mode=abs
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
  local container_id="${1:-}"
  local port="${2:-8080}"
  local env_file="${3:-.env}"

  if [ -z "$container_id" ]; then
    echo "Error: container image is required" >&2
    return 1
  fi
  if [[ ! "$port" =~ ^[0-9]+$ ]] || ((port < 1 || port > 65535)); then
    echo "Error: port must be between 1 and 65535" >&2
    return 1
  fi
  if [ ! -f "$env_file" ] || [ ! -r "$env_file" ]; then
    echo "Error: env file must be a readable file: $env_file" >&2
    return 1
  fi

  docker run \
    -v ~/.config/gcloud/application_default_credentials.json:/gcp/creds.json \
    --publish "${port}:8080" \
    --env-file "$env_file" \
    --env "PORT=8080" \
    --env "GOOGLE_APPLICATION_CREDENTIALS=/gcp/creds.json" \
    --env "GOOGLE_CLOUD_PROJECT=${PROJECT_ID:-}" \
    --interactive \
    "$container_id"
}

docker-build-and-run() {
  local env_file="${1:-.env}"
  docker build --platform linux/arm64/v8 -t swim-gen-backend:latest .
  docker-run swim-gen-backend:latest 8080 "$env_file"
}

scrape() {
  local url="${1:-}"
  if [ -z "$url" ]; then
    echo "Error: URL parameter is required"
    echo "Usage: task scrape -- <url>"
    exit 1
  fi

  shift
  go run ./cmd/scrape --url "$url" "$@"
}

# Check if the provided argument matches any of the functions
if [ -n "${1:-}" ] && ! declare -f "$1" >/dev/null; then
  echo "Error: Unknown task '$1'"
  echo
  help # Show help if the task is not recognized
  exit 1
fi

# Run application if no argument is provided
"${@:-run}"
