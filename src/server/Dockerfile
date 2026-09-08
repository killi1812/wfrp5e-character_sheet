# Stage 1: Build Go Server binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG BUILD=prod
ARG VERSION=0.0.0
ARG COMMIT_HASH=n/a
ARG BUILD_TIMESTAMP=n/a
ARG PACKAGE="github.com/killi1812/wfrp5e-character_sheet"

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X '${PACKAGE}/app.Build=${BUILD}' -X '${PACKAGE}/app.Version=${VERSION}' -X '${PACKAGE}/app.CommitHash=${COMMIT_HASH}' -X '${PACKAGE}/app.BuildTimestamp=${BUILD_TIMESTAMP}'" \
  -o wfrp5ecs main.go

# Stage 2: Lightweight runtime image
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/wfrp5ecs .

ENV PORT=8080
ENV MONGO_CONN=mongodb://database:27017

EXPOSE 8080

ENTRYPOINT ["./wfrp5ecs"]
