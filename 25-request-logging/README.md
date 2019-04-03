# 25. Request Logging

- Add a middleware to log something for each request.
- Requires creating a struct with some request values and passing it down through context.
- Ensure web.Respond updates the value.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/check.go
Modified cmd/sales-api/internal/handlers/products.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/main.go
Added    internal/mid/logger.go
Modified internal/platform/web/errors.go
Modified internal/platform/web/response.go
Modified internal/platform/web/web.go
```