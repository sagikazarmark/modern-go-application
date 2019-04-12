# Build image
FROM golang:1.12.3-alpine AS builder

ENV GOFLAGS="-mod=readonly"

RUN apk add --update --no-cache ca-certificates make git curl mercurial bzr

RUN mkdir -p /workspace
WORKDIR /workspace

COPY go.* /workspace/
RUN go mod download

COPY . /workspace

RUN BINARY_NAME=app make build-release


# Final image
FROM alpine:3.9.3

RUN apk add --update --no-cache ca-certificates tzdata bash curl

COPY --from=builder /workspace/build/${build_target}/app /usr/local/bin

EXPOSE 8000 8001 10000
CMD ["app", "--instrumentation.addr", ":10000", "--app.httpAddr", ":8000", "--app.grpcAddr", ":8001"]
