# Stage 1: Build the Vue.js application
FROM node:24-alpine AS build-stage

WORKDIR /app

COPY package*.json ./
RUN apk add --no-cache build-base autoconf automake libtool nasm zlib-dev
RUN npm install

COPY . .
RUN npm run build

# Stage 2: Serve the application with Nginx
FROM nginx:stable-alpine AS production-stage

COPY nginx.conf /etc/nginx/templates/default.conf.template

COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE ${PORT}
