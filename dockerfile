FROM golang:1.26-alpine AS builder

WORKDIR /src/server/build

COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY src/ .

# Define build arguments
ARG BUILD=prod
ARG VERSION=0.0.0
ARG COMMIT_HASH=n/a
ARG BUILD_TIMESTAMP=n/a
ARG PACKAGE="github.com/killi1812/wfrp5e-character_sheet"

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X '${PACKAGE}/app.Build=${BUILD}' -X '${PACKAGE}/app.Version=${VERSION}' -X '${PACKAGE}/app.CommitHash=${COMMIT_HASH}' -X '${PACKAGE}/app.BuildTimestamp=${BUILD_TIMESTAMP}'" \
  -o cache-server main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/wfrp5ecs .

# Re-declare build arguments to make them available in runtime stage ENV
ARG BUILD=prod
ARG VERSION=0.0.0
ARG COMMIT_HASH=n/a
ARG BUILD_TIMESTAMP=n/a

ENV APP_BUILD=${BUILD}
ENV APP_VERSION=${VERSION}
ENV APP_COMMIT_HASH=${COMMIT_HASH}
ENV APP_BUILD_TIMESTAMP=${BUILD_TIMESTAMP}

# Define entrypoint
ENTRYPOINT ["./wfrp5ecs"]
