# Loadgen

**A simple load generator tool for Modern Go Application**

Useful to demonstrate behavior and to fill charts with data.

```bash
docker build -t loadgen .
docker run --rm -it -e FRONTEND_ADDR=http://host.docker.internal:8000 loadgen
```

Alternatively, you can use the prebuilt image:

```bash
docker run --rm -it -e FRONTEND_ADDR=http://host.docker.internal:8000 sagikazarmark/modern-go-application:latest-loadgen
```
