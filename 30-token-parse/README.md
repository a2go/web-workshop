# 30. Authentication

- Add a context key type and variable for `auth.Claims`
- In package `mid` create a type `Auth` with a method `Authenticate`.
  - This method is a middleware that ensures the request has a valid token.
- Modify the `Handle` method in the `web` package to accept optional route specific middleware.
- Add the `Authenticate` middleware to all routes except the health check and the token generator.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/tests/products_test.go
Added    internal/mid/auth.go
Modified internal/mid/mid.go
Modified internal/platform/auth/auth.go
Modified internal/platform/auth/claims.go
Modified internal/platform/tests/tests.go
Modified internal/platform/web/web.go
```

## Dependency Changes:

```
+ 	go.opencensus.io v0.19.2
```