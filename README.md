# Modern Go Application

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go#project-layout)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagikazarmark/modern-go-application?style=flat-square)](https://goreportcard.com/report/github.com/sagikazarmark/modern-go-application)
[![GolangCI](https://golangci.com/badges/github.com/sagikazarmark/modern-go-application.svg)](https://golangci.com/r/github.com/sagikazarmark/modern-go-application)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/sagikazarmark/modern-go-application)

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sagikazarmark/modern-go-application/CI?style=flat-square)
[![CircleCI](https://circleci.com/gh/sagikazarmark/modern-go-application.svg?style=svg)](https://circleci.com/gh/sagikazarmark/modern-go-application)
[![Gitlab](https://img.shields.io/badge/gitlab-sagikazarmark%2Fmodern--go--application-orange.svg?logo=gitlab&longCache=true&style=flat-square)](https://gitlab.com/sagikazarmark/modern-go-application)

**Go application boilerplate and example applying modern practices**

This repository tries to collect the best practices of application development using Go language.
In addition to the language specific details, it also implements various language independent practices.

Some of the areas Modern Go Application touches:

- architecture
- package structure
- building the application
- testing
- configuration
- running the application (eg. in Docker)
- developer environment/experience
- telemetry

To help adopting these practices, this repository also serves as a boilerplate for new applications.


## Features

- configuration (using [spf13/viper](https://github.com/spf13/viper))
- logging (using [logur.dev/logur](https://logur.dev/logur) and [sirupsen/logrus](https://github.com/sirupsen/logrus))
- error handling (using [emperror.dev/emperror](https://emperror.dev/emperror))
- metrics and tracing using [Prometheus](https://prometheus.io/) and [Jaeger](https://www.jaegertracing.io/) (via [OpenCensus](https://opencensus.io/))
- health checks (using [AppsFlyer/go-sundheit](https://github.com/AppsFlyer/go-sundheit))
- graceful restart (using [cloudflare/tableflip](https://github.com/cloudflare/tableflip)) and shutdown
- support for multiple server/daemon instances (using [oklog/run](https://github.com/oklog/run))
- messaging (using [ThreeDotsLabs/watermill](https://github.com/ThreeDotsLabs/watermill))
- MySQL database connection (using [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
- ~~Redis connection (using [gomodule/redigo](https://github.com/gomodule/redigo))~~ removed due to lack of usage (see [#120](../../issues/120))


## First steps

To create a new application from the boilerplate clone this repository (if you haven't done already) into your GOPATH
then execute the following:

```bash
chmod +x init.sh && ./init.sh
? Package name (github.com/sagikazarmark/modern-go-application)
? Project name (modern-go-application)
? Binary name (modern-go-application)
? Service name (modern-go-application)
? Friendly service name (Modern Go Application)
? Update README (Y/n)
? Remove init script (y/N) y
```

It updates every import path and name in the repository to your project's values.
**Review** and commit the changes.


### Load generation

To test or demonstrate the application it comes with a simple load generation tool.
You can use it to test the example endpoints and generate some load (for example in order to fill dashboards with data).

Follow the instructions in [etc/loadgen](etc/loadgen).


## Inspiration

See [INSPIRATION.md](INSPIRATION.md) for links to articles, projects, code examples that somehow inspired
me while working on this project.


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
