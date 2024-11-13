ARG GOVERSION=1.22
FROM golang:${GOVERSION}-alpine

# add system requirements
RUN apk add --no-cache \
        ca-certificates \
        curl \
        findutils \
        just \
        git \
        tar \
        zip

# add ci tools

# golangci-lint
ARG golangci_lint_VERSION=v1.62.0
RUN curl -sSfL https://github.com/golangci/golangci-lint/releases/download/${golangci_lint_VERSION}/golangci-lint-${golangci_lint_VERSION#v}-linux-amd64.tar.gz \
        -o /tmp/golangci-lint.tar.gz \
 && tar -zxof /tmp/golangci-lint.tar.gz \
        --strip-components=1 \
        -C /go/bin/ \
        golangci-lint-${golangci_lint_VERSION#v}-linux-amd64/golangci-lint

# gotestsum
ARG gotestsum_VERSION=v1.12.0
RUN curl -sSfL https://github.com/gotestyourself/gotestsum/releases/download/${gotestsum_VERSION}/gotestsum_${gotestsum_VERSION#v}_linux_amd64.tar.gz \
        -o /tmp/gotestsum.tar.gz \
 && tar -zxof /tmp/gotestsum.tar.gz \
        -C /go/bin/ \
        gotestsum

