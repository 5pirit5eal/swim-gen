# BFF (Backend-for-Frontend) Service

This service acts as a Backend-for-Frontend for the Swim Gen Vue.js application. It sits between the static frontend and the Go backend API.

## Purpose

The primary responsibilities of this service are:

1. **Proxying Requests**: It forwards API requests from the frontend to the appropriate backend service.
2. **Authentication**: It securely handles service-to-service authentication. It injects a Google Cloud identity token into requests sent to the Go backend, ensuring that the backend is only accessible by this BFF service and not directly from the public internet.
3. **Simplifying the Frontend**: It abstracts away the details of backend service discovery and authentication from the client-side application.

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

    The server will be available at `http://localhost:8080`.

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
