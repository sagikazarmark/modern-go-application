# Build image
FROM golang:1.13.0-alpine AS builder

ENV GOFLAGS="-mod=readonly"

RUN apk add --update --no-cache ca-certificates make git curl mercurial bzr

RUN mkdir -p /workspace
WORKDIR /workspace

ARG GOPROXY

COPY go.* /workspace/
RUN go mod download

COPY . /workspace

ARG BUILD_TARGET

RUN set -xe && \
    if [[ "${BUILD_TARGET}" == "debug" ]]; then \
        cd /tmp; GOBIN=/workspace/build/debug go get github.com/go-delve/delve/cmd/dlv; cd -; \
        make build-debug; \
        mv build/debug /build; \
    else \
        make build-release; \
        mv build/release /build; \
    fi


# Final image
FROM alpine:3.10.1

RUN apk add --update --no-cache ca-certificates tzdata bash curl

SHELL ["/bin/bash", "-c"]

# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

ARG BUILD_TARGET

RUN if [[ "${BUILD_TARGET}" == "debug" ]]; then apk add --update --no-cache libc6-compat; fi

COPY --from=builder /build/* /usr/local/bin/

EXPOSE 8000 8001 10000
CMD ["modern-go-application", "--telemetry-addr", ":10000", "--http-addr", ":8000", "--grpc-addr", ":8001"]
