# Swim RAG Backend

This directory contains the backend service for the Swim RAG application. It is a Go-based RESTful API responsible for the core business logic, including training plan generation, database interactions, and PDF exports.

## Core Features

- **AI-Powered Plan Generation**: Leverages a Retrieval-Augmented Generation (RAG) system to create or recommend swimming training plans based on natural language queries.
- **Plan Donation**: Allows users to contribute new training plans to the system's database.
- **PDF Export**: Generates a PDF version of a training plan and uploads it to Google Cloud Storage.
- **Web Scraping**: Includes functionality to scrape training plans from external websites to populate the database.

## API Endpoints

The service exposes the following primary endpoints:

- `POST /query`: Queries the RAG system for a training plan.
- `POST /add`: Adds a new training plan to the database.
- `POST /export-pdf`: Exports a training plan to a PDF file.
- `GET /scrape`: Triggers the web scraping process.
- `POST /prompt`: Generates a prompt for the LLM.
- `GET /health`: Health check endpoint.
- `GET /swagger/*`: Serves the Swagger API documentation.

## Getting Started

### Prerequisites

- Go (version 1.22 or higher)
- Docker
- Google Cloud SDK

### Configuration

The application requires a `.env` file in the `backend` directory with the following environment variables:

```
PROJECT_ID=<your-gcp-project-id>
REGION=<your-gcp-region>
DB_USER=<your-database-user>
DB_NAME=<your-database-name>
DB_HOST=<your-database-host>
DB_PORT=<your-database-port>
DB_PASS_SECRET_ID=<your-secret-manager-secret-id>
GEMINI_API_KEY_SECRET_ID=<your-gemini-secret-id>
GCS_BUCKET_NAME=<your-gcs-bucket-name>
GCS_SERVICE_ACCOUNT_SECRET_ID=<your-gcs-sa-secret-id>
LOG_LEVEL=<log-level, e.g., DEBUG, INFO>
```

## Build, Test, and Run

This project uses a `Taskfile.sh` script to manage common tasks.

- **Run the application**:

    ```bash
    ./Taskfile.sh run
    ```

- **Run tests**:

    ```bash
    ./Taskfile.sh test
    ```

- **Format the code**:

    ```bash
    ./Taskfile.sh format
    ```

- **Validate the code (linting and vetting)**:

    ```bash
    ./Taskfile.sh validate
    ```

- **Generate Swagger docs**:

    ```bash
    ./Taskfile.sh docs
    ```

- **Build and run with Docker**:

    ```bash
    ./Taskfile.sh docker-build-and-run
    ```

## Key Technologies

- **Web Framework**: `chi`
- **Database**: `pgx` for PostgreSQL
- **AI/LLM**: `langchaingo` and `google.golang.org/genai`
- **PDF Generation**: `maroto`
- **Web Scraping**: `colly`
- **API Documentation**: `swaggo`
