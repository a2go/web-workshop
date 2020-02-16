# 28. Create Users

- Add a `internal/user` package with types `User` and `NewUser`.
- Add a `Create` function to the user package.
  - It should generate a password hash using bcrypt.
- Add a `internal/plaform/auth` package with constants `RoleAdmin` and `RoleUser`.
- Add a command `useradd` to the `sales-admin` program.


## File Changes:

```
Modified cmd/sales-admin/main.go
Added    internal/platform/auth/roles.go
Modified internal/schema/migrate.go
Modified internal/schema/seed.go
Added    internal/user/models.go
Added    internal/user/user.go
```

## Dependency Changes:

```
+ 	golang.org/x/crypto v0.0.0-20190320223903-b7391e95e576
```