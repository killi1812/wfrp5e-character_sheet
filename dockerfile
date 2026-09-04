# Stage 1: Build Vuetify Frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

COPY src/frontend/package*.json ./
RUN npm ci || npm install

COPY src/frontend ./
RUN npm run build

# Stage 2: Build Go Server
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/server

COPY src/server/go.mod src/server/go.sum ./
RUN go mod download

COPY src/server ./

ARG BUILD=prod
ARG VERSION=0.0.0
ARG COMMIT_HASH=n/a
ARG BUILD_TIMESTAMP=n/a
ARG PACKAGE="github.com/killi1812/wfrp5e-character_sheet"

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X '${PACKAGE}/app.Build=${BUILD}' -X '${PACKAGE}/app.Version=${VERSION}' -X '${PACKAGE}/app.CommitHash=${COMMIT_HASH}' -X '${PACKAGE}/app.BuildTimestamp=${BUILD_TIMESTAMP}'" \
  -o wfrp5ecs main.go

# Stage 3: Runtime Container (Single Container for Backend & Frontend)
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy backend binary
COPY --from=backend-builder /app/server/wfrp5ecs .

# Copy built static frontend assets to ./public
COPY --from=frontend-builder /app/frontend/dist ./public

ARG BUILD=prod
ARG VERSION=0.0.0
ARG COMMIT_HASH=n/a
ARG BUILD_TIMESTAMP=n/a

ENV APP_BUILD=${BUILD}
ENV APP_VERSION=${VERSION}
ENV APP_COMMIT_HASH=${COMMIT_HASH}
ENV APP_BUILD_TIMESTAMP=${BUILD_TIMESTAMP}
ENV PORT=8080
ENV MONGO_CONN=mongodb://mongo:27017

EXPOSE 8080

ENTRYPOINT ["./wfrp5ecs"]
