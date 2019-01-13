ARG GO_VERSION=1.11.4

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --update --no-cache ca-certificates make git curl mercurial

ARG PACKAGE=github.com/sagikazarmark/modern-go-application

RUN mkdir -p /go/src/${PACKAGE}
WORKDIR /go/src/${PACKAGE}

COPY Gopkg.* Makefile /go/src/${PACKAGE}/
RUN make vendor

COPY . /go/src/${PACKAGE}
RUN BUILD_DIR= BINARY_NAME=app make build-release


FROM alpine:3.8

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /app

EXPOSE 8000 10000
CMD ["/app", "--instrumentation.addr", ":10000", "--app.addr", ":8000"]
