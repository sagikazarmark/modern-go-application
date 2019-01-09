# Modern Go Application

[![Go Report Card](https://goreportcard.com/badge/github.com/sagikazarmark/modern-go-application?style=flat-square)](https://goreportcard.com/report/github.com/sagikazarmark/modern-go-application)
[![GolangCI](https://golangci.com/badges/github.com/sagikazarmark/modern-go-application.svg)](https://golangci.com/r/github.com/sagikazarmark/modern-go-application)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/sagikazarmark/modern-go-application)

[![Build Status](https://img.shields.io/travis/com/sagikazarmark/modern-go-application.svg?style=flat-square)](https://travis-ci.com/sagikazarmark/modern-go-application)
[![CircleCI (all branches)](https://img.shields.io/circleci/project/github/sagikazarmark/modern-go-application.svg?style=flat-square)](https://circleci.com/gh/sagikazarmark/modern-go-application)
[![Gitlab](https://img.shields.io/badge/gitlab-sagikazarmark%2Fmodern--go--application-orange.svg?logo=gitlab&longCache=true&style=flat-square)](https://gitlab.com/sagikazarmark/modern-go-application)

**Go application boilerplate and example applying modern practices**

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

- graceful restart (using [cloudflare/tableflip](https://github.com/cloudflare/tableflip)) and shutdown
- support for multiple server/daemon instances (using [oklog/run](https://github.com/oklog/run))
- metrics and tracing using [Prometheus](https://prometheus.io/) and [Jaeger](https://www.jaegertracing.io/) (via [OpenCensus](https://opencensus.io/))
- logging (using [goph/logur](https://github.com/goph/logur) and [sirupsen/logrus](https://github.com/goph/logur))
- health checks (using [InVisionApp/go-health](https://github.com/InVisionApp/go-health))
- configuration (using [spf13/viper](https://github.com/spf13/viper))


## Usage

To create a new application from the boilerplate clone this repository (if you haven't done already) and execute the following:

```bash
cd cmd/init && GO111MODULE=on go run . && cd -
? Package name github.com/sagikazarmark/modern-go-application
? Project name modern-go-application
? Binary name modern-go-application
? Remove this init script? Yes
```


## Inspiration

See [INSPIRATION.md](INSPIRATION.md) for links to articles, projects, code examples that somehow inspired
me while working on this project.


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
