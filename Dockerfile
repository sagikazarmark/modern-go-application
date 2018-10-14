FROM alpine:3.7

ARG BUILD_DIR
ARG BINARY_NAME

COPY $BUILD_DIR/$BINARY_NAME /service

EXPOSE 8000 10000
CMD ["/service", "--instrument-addr", ":10000", "--http-addr", ":8000"]
