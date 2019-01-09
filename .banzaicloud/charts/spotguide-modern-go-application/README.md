# Modern Go Application Helm Chart

## TL;DR

```bash
helm install -f values.local.yaml .
```


## Introduction

This chart bootstraps a Go application deployment on a
[Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.


## Prerequisites

- Kubernetes 1.10+


## Installing the Chart

To install the chart with the release name `my-release`:

```bash
helm install -f values.local.yaml --name my-release .
```

The command deploys the application on the Kubernetes cluster with the default configuration.
The configuration section lists the parameters that can be configured during installation.

> Tip: List all releases using `helm list`


## Configuration

The configurable parameters and default values are listed in [`values.yaml`](values.yaml).

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided during the chart installation:

```bash
helm install --name my-release -f my-values.yaml .
```
