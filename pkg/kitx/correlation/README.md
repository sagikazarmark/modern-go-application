# package correlation

**Package `correlation` provides a set of tools to add correlation ID to the context at certain levels (transport, endpoint) of the application.**

*Note:* Currently only server side middleware are implemented.

## Usage

### Get correlation ID from transport headers

Transport level middleware (`HTTPToContext`, `GRPCToContext`) read the correlation ID from headers
and add it to the context (if there is any).

```go
// HTTP example
httptransport.NewServer(
    ctx,
    endpoint,
    decoder,
    encoder,
    httptransport.ServerBefore(correlation.HTTPToContext()),
)

// gRPC example
grpctransport.NewServer(
    ctx,
    endpoint,
    decoder,
    encoder,
    grpctransport.ServerBefore(correlation.GRPCToContext()),
)
```

### Generate a correlation ID if none is found in the context

When clients don't pass a correlation ID to the server, one should be generated early of the request lifecycle.

```go
endpoint = correlation.Middleware()(endpoint)
```

Make sure to seed the global random number generator:

```go
import (
    "math/rand"
    "time"
)

rand.Seed(time.Now().UnixNano())
```
