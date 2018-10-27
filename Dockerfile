FROM golang:1.11-alpine AS builder

RUN apk add --update --no-cache ca-certificates make git curl mercurial

ARG PACKAGE=github.com/sagikazarmark/modern-go-application

RUN mkdir -p /go/src/${PACKAGE}
WORKDIR /go/src/${PACKAGE}

COPY Gopkg.* Makefile /go/src/${PACKAGE}/
RUN make vendor

COPY . /go/src/${PACKAGE}
RUN BUILD_DIR=/tmp BINARY_NAME=service make build-release


FROM alpine:3.7

COPY --from=builder /tmp/service /service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000 10000
CMD ["/service", "--maintenance-addr", ":10000", "--http-addr", ":8000"]
