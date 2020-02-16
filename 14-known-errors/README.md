# 14. Known Errors

This error pattern effectively treats all errors as fatal:

```go
if err != nil {
	return err
}
```

Document expected errors in the `products` package and return them when they
occur. The handlers should look for these errors and use `web.NewRequestError`
to send them up the call chain with an HTTP status.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified internal/product/product.go
```