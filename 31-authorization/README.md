# 31. Authorization

- Add a method to `Claims` called `HasRole`.
- Add a method `HasRole` to the `Auth` middleware type.
- Add the `HasRole` middleware requiring `RoleAdmin` on the Delete and AddSale routes.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/routes.go
Modified internal/mid/auth.go
Modified internal/platform/auth/claims.go
```