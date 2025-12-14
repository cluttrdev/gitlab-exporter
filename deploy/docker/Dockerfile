# syntax=docker/dockerfile:1

# ============================================================================
# Generate CA certificates
FROM alpine:3.23.0 AS certs

RUN apk --no-cache add ca-certificates

# ============================================================================
# Shared base image
FROM scratch AS base

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY <<EOF /etc/passwd
nobody:x:65534:65534:nobody:/nonexistent:/bin/false
EOF

COPY <<EOF /etc/group
nobody:x:65534:
EOF

USER nobody:nobody

# ============================================================================
# Runtime Image Targets
# ============================================================================

# gitlab-exporter runtime image
FROM base AS gitlab-exporter
COPY --chown=nobody:nobody gitlab-exporter /bin/gitlab-exporter
ENTRYPOINT [ "/bin/gitlab-exporter" ]

# gitlab-exporter-clickhouse-recorder runtime image
FROM base AS gitlab-exporter-clickhouse-recorder
COPY --chown=nobody:nobody gitlab-exporter-clickhouse-recorder /bin/gitlab-exporter-clickhouse-recorder
ENTRYPOINT [ "/bin/gitlab-exporter-clickhouse-recorder" ]
