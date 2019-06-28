# 36. Self Shutdown

Certain error conditions only occur because of programmer error. In those known
cases we can't let the service continue to run. The application must shut down.

- Pass the signal channel from `func main` down to the web framework.
- Add a special error type in `web/errors.go` that represents a shutdown error.
- Let the error return up from the `ErrorHandler` middleware.
- Make the top level web function detect unhandled errors and use the channel to shut the application down.

## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/internal/handlers/user.go
Modified cmd/sales-api/main.go
Modified cmd/sales-api/tests/product_test.go
Modified cmd/sales-api/tests/user_test.go
Modified internal/mid/errors.go
Modified internal/mid/logger.go
Modified internal/mid/panics.go
Modified internal/platform/web/errors.go
Modified internal/platform/web/response.go
Modified internal/platform/web/web.go
Modified internal/product/product.go
Modified internal/product/product_test.go
Modified internal/product/sales.go
```