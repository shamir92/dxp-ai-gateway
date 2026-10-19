# Build stage: compile a static binary.
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies separately from source.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off keeps the binary static so it can run on a bare Alpine base.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Runtime stage.
FROM alpine:3.23

# ca-certificates for outbound HTTPS (external platforms, OAuth, SMTP, geoip).
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app

WORKDIR /app

COPY --from=build /out/server /app/server

# These two are read from the filesystem at startup, relative to the working
# directory, so they must exist in the image:
#   migrations                       -> runMigrations (file://migrations)
#   internal/pkg/mail/templates      -> MAIL_TEMPLATES_DIR default
COPY --from=build /src/migrations /app/migrations
COPY --from=build /src/internal/pkg/mail/templates /app/internal/pkg/mail/templates

USER app

ENV HTTP_ADDR=:8080 \
    MAIL_TEMPLATES_DIR=internal/pkg/mail/templates

EXPOSE 8080

# /healthz is unauthenticated and cheap, so it suits both probe kinds.
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -q -O- http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/app/server"]