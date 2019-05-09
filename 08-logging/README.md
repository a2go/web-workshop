# 8. Logging

- Do not use any package level variables such as the global `log`.
- Pass a `*log.Logger` to dependencies.

## Notes

Log actionable events. This is separate from Metrics or Tracing.

## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/main.go
```