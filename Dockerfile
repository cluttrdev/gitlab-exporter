# syntax=docker/dockerfile:1

# ============================================================================
# Shared base image
FROM scratch AS base

COPY --from=docker.io/alpine:3.23.0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY <<EOF /etc/passwd
nobody:x:65534:65534:nobody:/nonexistent:/bin/false
EOF

COPY <<EOF /etc/group
nobody:x:65534:
EOF

USER nobody:nobody

# ============================================================================
# Runtime Image
# ============================================================================
FROM base
ARG APP
ARG TARGETOS
ARG TARGETARCH
COPY --chown=nobody:nobody ${TARGETOS}_${TARGETARCH}/$APP /bin/app
ENTRYPOINT [ "/bin/app" ]
