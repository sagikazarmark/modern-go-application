FROM scratch

ARG BUILD_DIR
ARG BINARY_NAME

COPY $BUILD_DIR/$BINARY_NAME /service

EXPOSE 10000
CMD ["/service", "--debug-addr", ":10000"]
