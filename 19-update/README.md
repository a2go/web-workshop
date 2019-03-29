# 19. Update

- Add `Update` function in `internal/products`.
- Use a defined type `UpdateProduct` as an argument.
- To support partial updates, allow the fields of this type to be null.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/products.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/tests/products_test.go
Added    internal/platform/tests/tests.go
Modified internal/platform/web/response.go
Modified internal/products/products.go
Modified internal/products/products_test.go
```