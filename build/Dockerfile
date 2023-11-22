# syntax=docker/dockerfile:1

#
# BUILD
#

FROM docker.io/golang:1.20-alpine3.18 AS build

RUN apk add --update --no-cache ca-certificates

RUN adduser \
    --disabled-password \
    --shell /sbin/nologin \
    --home /nonexistent \
    --no-create-home \
    --gecos "" \
    glche

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download

COPY . /build/

ARG GOOS=linux
ARG GOARCH=amd64
ARG VERSION=v0.0.0+unknown

RUN GOOS=${GOOS} GOARCH=${GOARCH} go build \
    -o /bin/gitlab-exporter \
    -ldflags "-X 'main.version=${VERSION}' -s -w"

#
# RUN
#

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /bin/gitlab-exporter /bin/gitlab-exporter

USER glche:glche

ENTRYPOINT [ "/bin/gitlab-exporter" ]

CMD [ "--help" ]
