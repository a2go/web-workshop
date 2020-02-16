# 26. Request Logging

- Add a middleware to log something for each request.
- Requires creating a struct with some request values and passing it down through context.
- Ensure web.Respond updates the value.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/check.go
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified internal/mid/errors.go
Added    internal/mid/logger.go
Modified internal/platform/web/response.go
Modified internal/platform/web/web.go
```