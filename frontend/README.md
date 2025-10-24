# Swim-RAG Frontend

Vue 3 + TypeScript frontend for the Swim Training Plan Generator application.

The app can be displayed with both mobile and larger displays. The latter is recommended.

## Key Technologies

- **Vue 3**: The progressive JavaScript framework for building user interfaces.
- **TypeScript**: A typed superset of JavaScript that compiles to plain JavaScript.
- **Vite**: A fast build tool that provides a quicker and leaner development experience.
- **Pinia**: The official state management library for Vue.
- **Vue Router**: The official router for Vue.js.
- **Vitest**: A blazing fast unit-test framework powered by Vite.
- **ESLint**: Pluggable and configurable linter tool for identifying and reporting on patterns in JavaScript.
- **Prettier**: An opinionated code formatter.

## Project Setup

**Prerequisites:**

- Node.js >=22.0.0

```sh
npm install
```

### Development Commands

```sh
npm run dev             # Start development server with HMR
npm run build           # Type-check and build for production
npm run build-only      # Build for production without type-checking
npm run type-check      # Run TypeScript validation
npm run test:unit       # Run unit tests once
npm run test:continous  # Run unit tests in watch mode
npm run lint            # ESLint with auto-fix
npm run format          # Prettier formatting
npm run preview         # Preview production build locally
npm run optimize-images # Optimize images for production
```

### Running with the BFF and Backend

To run the frontend with the BFF locally, you need to have both the backend and the BFF services running.

1. **Start the backend service** by following the instructions in the `backend/README.md` file. By default, it runs on port `8080`.
2. **Start the BFF service** by following the instructions in the `bff/README.md` file. By default, it runs on port `8081`.
3. **Start the frontend development server**:

    ```sh
    npm run dev
    ```

    The `vite.config.ts` is configured to proxy requests from `/api` to the BFF service running on `http://localhost:8081`.

### Running with Docker Compose

To run the entire application stack (frontend, BFF, and backend) using Docker, you can use the `docker-compose.yml` file located in the root of the project.

1. **Ensure you have Docker and Docker Compose installed.**
2. **Navigate to the root of the `swim-gen` project.**
3. **Run the following command:**

    ```sh
    docker-compose up --build
    ```

    This will build the Docker images for the frontend, BFF, and backend services and start them. The frontend will be available at `http://localhost:5173`.
