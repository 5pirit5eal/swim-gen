#!/bin/bash

# Load testing script for the RAG service on Cloud Run

# --- Configuration ---
SERVICE_URL=http://localhost:8080
NUM_USERS=1000 # Number of concurrent users to simulate
REQUEST_DELAY=0.1 # Delay between requests in seconds

# --- Counters ---
SUCCESS_COUNT=0
ERROR_COUNT=0

# Create temporary files for counters
SUCCESS_FILE=$(mktemp)
ERROR_FILE=$(mktemp)
echo 0 > $SUCCESS_FILE
echo 0 > $ERROR_FILE

# --- Functions ---

# Function to make a request to the /generate-prompt endpoint
generate_prompt() {
    local lang=$1
    local token=$(gcloud auth print-identity-token)

    curl -s -X POST "$SERVICE_URL/prompt" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{\"language\": \"$lang\"}"
}

# Function to make a request to the /query endpoint
query_plan() {
    local prompt=$1
    local token=$(gcloud auth print-identity-token)

    curl -s -X POST "$SERVICE_URL/query" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{\"content\": \"$prompt\", \"method\": \"generate\"}"
}

# --- Main Execution ---

echo "Starting load test with $NUM_USERS concurrent users..."

for i in $(seq 1 $NUM_USERS); do
    (
        echo "User $i: Generating prompt..."
        PROMPT_RESPONSE=$(generate_prompt "en")

        # Check if the prompt generation was successful
        if ! echo "$PROMPT_RESPONSE" | jq -e .prompt > /dev/null; then
            echo "User $i: Failed to generate or parse prompt. Response: $PROMPT_RESPONSE"
            echo $(($(cat $ERROR_FILE) + 1)) > $ERROR_FILE
            exit 1
        fi

        PROMPT=$(echo "$PROMPT_RESPONSE" | jq -r .prompt)
        echo "User $i: Querying with prompt: $PROMPT"

        QUERY_RESPONSE=$(query_plan "$PROMPT")

        # Check if the query was successful
        if ! echo "$QUERY_RESPONSE" | jq -e .table > /dev/null; then
            echo "User $i: Failed to query plan. Response: $QUERY_RESPONSE"
            echo $(($(cat $ERROR_FILE) + 1)) > $ERROR_FILE
        else
            echo "User $i: Successfully queried plan"
            echo $(($(cat $SUCCESS_FILE) + 1)) > $SUCCESS_FILE
        fi
    ) &
    sleep $REQUEST_DELAY
done

wait

# --- Aggregate Results ---
SUCCESS_COUNT=$(cat $SUCCESS_FILE)
ERROR_COUNT=$(cat $ERROR_FILE)

# Clean up temporary files
rm $SUCCESS_FILE
rm $ERROR_FILE

# --- Results ---
echo "---------------------"
echo "Load Test Results"
echo "---------------------"
echo "Total Requests: $NUM_USERS"
echo "Successful Requests: $SUCCESS_COUNT"
echo "Failed Requests: $ERROR_COUNT"
echo "---------------------"

if [ $ERROR_COUNT -gt 0 ]; then
    exit 1
fi
