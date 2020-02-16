# 20. Update

- Add `Update` function in `internal/product`.
- Use a defined type `UpdateProduct` as an argument.
- To support partial updates, allow the fields of this type to be null.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/tests/product_test.go
Modified internal/platform/web/response.go
Modified internal/product/models.go
Modified internal/product/product.go
Modified internal/product/product_test.go
Modified internal/tests/tests.go
```