# Stage 1: Build the application
FROM node:22-slim AS builder

# Set the working directory
WORKDIR /usr/src/app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install all dependencies (including devDependencies)
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the TypeScript code
RUN npm run build

# Remove devDependencies after build
RUN npm prune --production

# Stage 2: Create the production image
FROM node:22-slim

# Set the working directory
WORKDIR /usr/src/app

# Copy the pruned node_modules and the built code from the builder stage
COPY --from=builder /usr/src/app/node_modules ./node_modules
COPY --from=builder /usr/src/app/dist ./dist

# Expose the port the app runs on
EXPOSE 8080

# Define the command to run the application
CMD [ "node", "dist/main.js" ]
