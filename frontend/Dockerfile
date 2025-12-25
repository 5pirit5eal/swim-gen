# Stage 1: Build the Vue.js application
FROM node:24-alpine AS build-stage

WORKDIR /app

ARG VITE_IMPRESSUM_NAME
ARG VITE_IMPRESSUM_ADDRESS
ARG VITE_IMPRESSUM_CITY
ARG VITE_IMPRESSUM_PHONE
ARG VITE_IMPRESSUM_EMAIL
ARG VITE_SITE_URL
ARG VITE_SUPABASE_URL
ARG VITE_SUPABASE_ANON_KEY
ARG VITE_PUBLIC_BUCKET_NAME

COPY package*.json ./
RUN apk add --no-cache build-base autoconf automake libtool nasm zlib-dev

# Set args as env variables for npm build process
ENV VITE_IMPRESSUM_NAME=$VITE_IMPRESSUM_NAME
ENV VITE_IMPRESSUM_ADDRESS=$VITE_IMPRESSUM_ADDRESS
ENV VITE_IMPRESSUM_CITY=$VITE_IMPRESSUM_CITY
ENV VITE_IMPRESSUM_PHONE=$VITE_IMPRESSUM_PHONE
ENV VITE_IMPRESSUM_EMAIL=$VITE_IMPRESSUM_EMAIL
ENV VITE_SITE_URL=$VITE_SITE_URL
ENV VITE_SUPABASE_URL=$VITE_SUPABASE_URL
ENV VITE_SUPABASE_ANON_KEY=$VITE_SUPABASE_ANON_KEY
ENV VITE_PUBLIC_BUCKET_NAME=$VITE_PUBLIC_BUCKET_NAME

RUN npm install

COPY . .

# Delete any env variable files
RUN rm -f .env .env.*

RUN npm run build

# Stage 2: Serve the application with Nginx
FROM nginx:stable-alpine AS production-stage

COPY nginx.conf /etc/nginx/templates/default.conf.template

COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE ${PORT}
