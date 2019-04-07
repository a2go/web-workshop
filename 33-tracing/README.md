# 33. Tracing

- Add a second service to `docker-compose.yaml` to run OpenZipkin. Start that service.
- Create a zipkin exporter in `cmd/sales-api/main.go` and register it with the OpenCensus trace package.
  - Define configuration parameters for the trace url, service name, and trace probability.
- Ensure the reporter is closed when `func run` returns so final values are flushed.
- Starting with the adapter function in `web.go`, begin the collection of spans. Most will be like this:

```go
ctx, span := trace.StartSpan(ctx, "handlers.Products.List")
defer span.End()
```

Generate some load then view the traces at http://localhost:9411/

## Notes:

Not every span has to match to a specific function. You can identify narrower
spans by explicitly using `StartSpan` and `span.End`. See `mid/auth.go`.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/check.go
Modified cmd/sales-api/internal/handlers/products.go
Modified cmd/sales-api/internal/handlers/users.go
Modified cmd/sales-api/main.go
Modified docker-compose.yaml
Modified internal/mid/auth.go
Modified internal/mid/errors.go
Modified internal/mid/logger.go
Modified internal/mid/metrics.go
Modified internal/platform/web/web.go
Modified internal/products/products.go
Modified internal/products/sales.go
Modified internal/users/users.go
```

## Dependency Changes:

```
+ 	github.com/openzipkin/zipkin-go v0.1.6
```