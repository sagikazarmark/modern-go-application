# Modern Go Application

This repository has multiple purposes:

- It serves as a boilerplate for new projects
- It tries to collect the best practices for various areas of developing a (modern) application

It tries to include many things related to application development:

- architecture
- package structure
- building the application
- testing
- configuration
- running the application (eg. Docker)
- developer environment/experience
- instrumentation


Some of the features:

- graceful reload (using [github.com/cloudflare/tableflip](https://github.com/cloudflare/tableflip)) and shutdown
- support for multiple server/daemon instances (using [github.com/oklog/run](https://github.com/oklog/run))
- metrics and tracing using [Prometheus](https://prometheus.io/) and [Jaeger](https://www.jaegertracing.io/) (via [OpenCensus](https://opencensus.io/))
- logging (using [github.com/go-kit/kit](https://github.com/go-kit/kit))
- health checks (using [github.com/InVisionApp/go-health](https://github.com/InVisionApp/go-health))
- configuration (using [github.com/spf13/viper](https://github.com/spf13/viper))


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
