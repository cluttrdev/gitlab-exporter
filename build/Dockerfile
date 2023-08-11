# synta=docker/dockerfile:1

#
# BUILD
#

FROM docker.io/golang:1.20-alpine3.18 AS build

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download

COPY . /build/

RUN go build -o /gitlab-clickhouse-exporter

#
# RUN
#

FROM docker.io/alpine:3.18

WORKDIR /

COPY --from=build /gitlab-clickhouse-exporter /gitlab-clickhouse-exporter

# run as nobody
USER 65534:65534

ENTRYPOINT [ "/gitlab-clickhouse-exporter" ]

CMD [ "--help"]
