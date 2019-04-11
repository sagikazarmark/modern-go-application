# Build target
ARG build_target=release

FROM alpine:3.9.3

ARG build_target

RUN if [[ "${build_target}" != "debug" && "${build_target}" != "release" ]]; then echo -e "\033[0;31mBuild argument \$build_target must be \"release\" or \"debug\".\033[0m" && exit 1; fi


# Base run image
FROM alpine:3.9.3 AS base-release

RUN apk add --update --no-cache ca-certificates tzdata bash curl

EXPOSE 8000 8001 10000
CMD ["app", "--instrumentation.addr", ":10000", "--app.httpAddr", ":8000", "--app.grpcAddr", ":8001"]


# Build image
FROM golang:1.12.3-alpine AS builder

ENV GOFLAGS="-mod=readonly"

RUN apk add --update --no-cache ca-certificates make git curl mercurial bzr

RUN mkdir -p /workspace
WORKDIR /workspace

COPY go.* /workspace/
RUN go mod download

COPY . /workspace

ARG build_target

RUN set -xe; if [[ "${build_target}" == "debug" ]]; then cd /tmp; go get github.com/derekparker/delve/cmd/dlv; cd -; else touch /go/bin/dlv; fi

RUN BINARY_NAME=app make build-${build_target}


# Base debug image
FROM base-release AS base-debug

COPY --from=builder /go/bin/dlv /usr/local/bin

EXPOSE 8000 8001 10000 40000
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--log", "exec", "/usr/local/bin/app", "--", "--instrumentation.addr", ":10000", "--app.httpAddr", ":8000", "--app.grpcAddr", ":8001"]


# Final image
FROM base-$build_target

ARG build_target

COPY --from=builder /workspace/build/${build_target}/app /usr/local/bin
