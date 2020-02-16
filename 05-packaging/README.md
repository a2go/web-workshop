# 5. Packaging

- Put business logic for Products in `internal/products`
- Put db administration in `cmd/sales-admin`
- Put entrypoint in `cmd/sales-api`
- Put HTTP layer in `cmd/sales-api/internal/handlers`

## Links:

https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

## File Changes:

```
Added    cmd/sales-admin/main.go
Added    cmd/sales-api/internal/handlers/product.go
Added    cmd/sales-api/main.go
Added    internal/platform/database/database.go
Added    internal/product/models.go
Added    internal/product/product.go
Moved    schema/migrate.go -> internal/schema/migrate.go
Moved    schema/seed.go -> internal/schema/seed.go
Deleted  main.go
```

## Dependency Changes:

```
+ 	github.com/pkg/errors v0.8.1
```