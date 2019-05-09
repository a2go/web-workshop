# 29. Login

- Add a `keygen` command to `sales-admin` which can generate a x509 private key.
- Add a `Claims` type to the `auth` package.
- Add an `Authenticator` type to the `auth` package which generates a JWT for `Claims`.
- Add an `Authenticate` function to the `users` package which finds a user by email and verifies their password.
- Add a handler that identifies a user from Basic Auth and responds with a Token.


## File Changes:

```
Modified cmd/sales-admin/main.go
Modified cmd/sales-api/internal/handlers/routes.go
Added    cmd/sales-api/internal/handlers/user.go
Modified cmd/sales-api/main.go
Modified cmd/sales-api/tests/product_test.go
Added    cmd/sales-api/tests/user_test.go
Added    internal/platform/auth/auth.go
Modified internal/platform/auth/claims.go
Modified internal/platform/tests/tests.go
Modified internal/user/user.go
Added    private.pem
```

## Dependency Changes:

```
+ 	github.com/dgrijalva/jwt-go v3.2.0+incompatible
```