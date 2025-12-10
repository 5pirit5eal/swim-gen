# BFF (Backend-for-Frontend) Service

This service acts as a Backend-for-Frontend for the Swim Gen Frontend application. It sits between the static frontend and the Go backend API.

## Purpose

The primary responsibilities of this service are:

1. **Proxying Requests**: It forwards API requests from the frontend to the appropriate backend service.
2. **Dual Authentication Handling**:
   - **User Authentication**: Passes through the user's Supabase access token from the frontend's `Authorization` header to the backend for user identification and authorization.
   - **Service-to-Service Authentication**: Adds a Google Cloud Identity token (via `X-Serverless-Authorization` header) for Cloud Run service-to-service authentication, ensuring the backend is only accessible by this BFF service and not directly from the public internet.
3. **Rate Limiting**: Protects the backend by applying rate limiting to API requests.
4. **CORS Management**: Handles Cross-Origin Resource Sharing (CORS) configuration to allow only authorized frontend origins.
5. **Simplifying the Frontend**: Abstracts away the details of backend service discovery and authentication from the client-side application.

## Authentication Flow

### User Authentication (Supabase)

- Frontend sends user's Supabase access token in the `Authorization: Bearer <token>` header
- BFF passes this header through to the backend unchanged
- Backend verifies the Supabase token and extracts user identity
- Anonymous requests (no Authorization header) are also supported

### Service-to-Service Authentication (Google Cloud Run)

- When deployed to Cloud Run, the BFF adds a Google Identity token to the `X-Serverless-Authorization` header
- This token is obtained using the `google-auth-library` package
- Cloud Run verifies this token to authenticate the BFF service
- In development mode (`NODE_ENV=development`), this token is skipped for local testing

## Logging Configuration

The BFF uses a configurable logging system controlled by the `LOG_LEVEL` environment variable.

### Supported Log Levels

| Level   | Description                                      |
|---------|--------------------------------------------------|
| `DEBUG` | Detailed debugging info (e.g., request proxying) |
| `INFO`  | General informational messages (default)         |
| `WARN`  | Warning messages                                 |
| `ERROR` | Error messages only                              |

### Configuration

Set the `LOG_LEVEL` environment variable to control which messages are logged:

```env
LOG_LEVEL=DEBUG  # Enable all logging including request details
LOG_LEVEL=INFO   # Default - standard informational messages
LOG_LEVEL=WARN   # Only warnings and errors
LOG_LEVEL=ERROR  # Only errors
```

When deployed, this is configured via the `log_level` Terraform variable in `deployments/prod/1-services/terraform.tfvars` or `deployments/dev/1-services/terraform.tfvars`.

## Tech Stack

- **Runtime**: Node.js
- **Framework**: Express.js
- **Language**: TypeScript
- **HTTP Client**: axios
- **Authentication**: google-auth-library

## Development

### Prerequisites

- Node.js (v18 or higher)
- npm

### Installation

1. Navigate to the `bff` directory:

    ```bash
    cd bff
    ```

2. Install dependencies:

    ```bash
    npm install
    ```

### Running in Development Mode

1. Create a `.env.development` file in the `bff` directory. This file is used to configure the service for local development and is not committed to source control.

    ```env
    # The full URL of the Go backend service when running locally.
    BACKEND_URL=http://localhost:8081
    ```

2. Start the development server, which uses `nodemon` to automatically restart on file changes:

    ```bash
    npm run dev
    ```

    The server will be available at `http://localhost:8081`.

## Building for Production

To compile the TypeScript code into JavaScript in the `dist` directory, run:

```bash
npm run build
```

To run the compiled code, use:

```bash
npm run start
```

## Testing

To run the test suite, use:

```bash
npm run test
```
