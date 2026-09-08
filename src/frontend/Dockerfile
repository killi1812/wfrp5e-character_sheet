# Stage 1: Build Vue 3 Frontend
FROM node:22-alpine AS builder
WORKDIR /app

COPY package*.json ./
RUN npm ci || npm install

COPY . ./
RUN npm run build

# Stage 2: Lightweight Caddy static SPA server
FROM caddy:alpine
COPY --from=builder /app/dist /usr/share/caddy

RUN printf ":80 {\n\troot * /usr/share/caddy\n\tfile_server\n\ttry_files {path} /index.html\n}\n" > /etc/caddy/Caddyfile

EXPOSE 80
