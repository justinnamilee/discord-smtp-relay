# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Install build deps (and ca certs for static binary TLS if ever needed)
RUN apk add --no-cache git ca-certificates

# Cache go modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build a static binary
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/discord-smtp-relay ./main.go

# ---- Runtime stage ----
# You can also use gcr.io/distroless/static:nonroot
FROM alpine:3.20
WORKDIR /app

# Certificates (HTTPS to Discord webhooks)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary
COPY --from=builder /out/discord-smtp-relay /app/discord-smtp-relay

# Copy example templates (optional; mount your own in prod)
# You can remove this if you only ever mount /app/etc at runtime.
COPY etc /app/etc

# Non-root user
RUN adduser -D -H -u 10001 appuser
USER appuser

# Default port in the app is 1025; make it explicit
EXPOSE 1025/tcp

# Healthcheck: TCP open on the SMTP port
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD /bin/sh -c "exec 3<>/dev/tcp/127.0.0.1/${PORT:-1025} && exec 3>&-"

# Entrypoint
ENTRYPOINT ["/app/discord-smtp-relay"]
