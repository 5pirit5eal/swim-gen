# Stage 1: Build the Vue.js application
FROM node:24-alpine AS build-stage

WORKDIR /app

ARG VITE_IMPRESSUM_NAME
ARG VITE_IMPRESSUM_ADDRESS
ARG VITE_IMPRESSUM_CITY
ARG VITE_IMPRESSUM_PHONE
ARG VITE_IMPRESSUM_EMAIL
ARG VITE_SITE_URL

COPY package*.json ./
RUN apk add --no-cache build-base autoconf automake libtool nasm zlib-dev
RUN npm install

COPY . .

RUN echo "VITE_IMPRESSUM_NAME=${VITE_IMPRESSUM_NAME}" > .env.production && \
    echo "VITE_IMPRESSUM_ADDRESS=${VITE_IMPRESSUM_ADDRESS}" >> .env.production && \
    echo "VITE_IMPRESSUM_CITY=${VITE_IMPRESSUM_CITY}" >> .env.production && \
    echo "VITE_IMPRESSUM_PHONE=${VITE_IMPRESSUM_PHONE}" >> .env.production && \
    echo "VITE_IMPRESSUM_EMAIL=${VITE_IMPRESSUM_EMAIL}" >> .env.production && \
    echo "VITE_SITE_URL=${VITE_SITE_URL}" >> .env.production

RUN npm run build

# Stage 2: Serve the application with Nginx
FROM nginx:stable-alpine AS production-stage

COPY nginx.conf /etc/nginx/templates/default.conf.template

COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE ${PORT}
