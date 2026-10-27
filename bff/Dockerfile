# Stage 1: Build the application
FROM node:26-slim@sha256:36793bcbb160a6a16cede1e62b93f84e26b78a972a8a4c0b879c1c12f22761ec AS builder

# Set the working directory
WORKDIR /usr/src/app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install all dependencies (including devDependencies)
RUN npm ci

# Copy the rest of the application code
COPY . .

# Build the TypeScript code
RUN npm run build

# Remove devDependencies after build
RUN npm prune --omit=dev

# Stage 2: Create the production image
FROM node:26-slim@sha256:36793bcbb160a6a16cede1e62b93f84e26b78a972a8a4c0b879c1c12f22761ec

# Set the working directory
WORKDIR /usr/src/app

# Copy the pruned node_modules and the built code from the builder stage
COPY --from=builder --chown=node:node /usr/src/app/node_modules ./node_modules
COPY --from=builder --chown=node:node /usr/src/app/dist ./dist

# Run the application as the built-in non-root Node user.
USER node

# Expose the port the app runs on
EXPOSE 8080

# Define the command to run the application
CMD [ "node", "dist/main.js" ]
