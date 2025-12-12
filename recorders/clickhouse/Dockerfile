# syntax=docker/dockerfile:1

#
# BUILD
#

ARG GOVERSION=1.24

FROM docker.io/golang:${GOVERSION}-alpine AS build

RUN apk add --update --no-cache ca-certificates

RUN adduser \
    --disabled-password \
    --shell /sbin/nologin \
    --home /nonexistent \
    --no-create-home \
    --gecos "" \
    user

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download

COPY . /build/

ARG GOOS=linux
ARG GOARCH=amd64
ARG VERSION=v0.0.0+unknown

RUN GOOS=${GOOS} GOARCH=${GOARCH} go build \
    -ldflags "-X 'main.version=${VERSION}' -s -w" \
    -o /bin/gitlab-exporter-clickhouse-recorder

#
# RUN
#

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /bin/gitlab-exporter-clickhouse-recorder /bin/gitlab-exporter-clickhouse-recorder

USER user:user

ENTRYPOINT [ "/bin/gitlab-exporter-clickhouse-recorder" ]

CMD [ "--help" ]
