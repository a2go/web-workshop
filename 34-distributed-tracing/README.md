# 34. Distributed Tracing

In a microservice environment it can be hard to see the execution path of a
single request across multiple services. The tracing libraries / products we
are using support propogating context across requests so they are correlated as
part of the same overall request.

Look for a Trace ID on incoming requests and add a TraceID to outgoing requests.


## File Changes:

```
Modified internal/mid/errors.go
Modified internal/mid/logger.go
Added    internal/platform/tracing/tracing.go
Modified internal/platform/web/web.go
```