# 25. Metrics

- Import `expvar` to expose custom variables to external clients.
- Add middleware to track number of requests, number of errors, and current goroutine count.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/main.go
Added    internal/mid/metrics.go
```