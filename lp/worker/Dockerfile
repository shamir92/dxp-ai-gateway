# Provisioning worker image. The build context is this directory, so the worker
# builds independently of the backend and frontend:
#   docker build -t dxp-ai-gateway/worker lp/worker
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies separately from source.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off keeps the binary static so it can run on a bare Alpine base.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/worker .

# Runtime stage. The worker serves no HTTP and reads no files from disk, so it
# only needs CA certificates for outbound HTTPS to the external platforms.
FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app

WORKDIR /app

COPY --from=build /out/worker /app/worker

USER app

# Release identity for Sentry: the git tag on a tagged build, otherwise the sha.
# Overridable at run time; an explicit SENTRY_RELEASE in the environment wins.
ARG SENTRY_RELEASE=""
ENV SENTRY_RELEASE=$SENTRY_RELEASE

ENTRYPOINT ["/app/worker"]
