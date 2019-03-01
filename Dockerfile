ARG GO_VERSION=1.12

FROM golang:${GO_VERSION}-alpine AS builder

ENV GOFLAGS="-mod=readonly"

RUN apk add --update --no-cache ca-certificates make git curl mercurial bzr

RUN mkdir -p /build
WORKDIR /build

COPY go.* /build/
RUN go mod download

COPY . /build
RUN BINARY_NAME=app make build-release


FROM alpine:3.9

RUN apk add --update --no-cache ca-certificates tzdata

COPY --from=builder /build/build/release/app /app

EXPOSE 8000 10000
CMD ["/app", "--instrumentation.addr", ":10000", "--app.addr", ":8000"]
