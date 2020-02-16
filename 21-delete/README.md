# 21. Delete

- Add `Delete` function in `internal/products`,
- Add route for `DELETE` `/v1/products/{id}`.
- Add `Delete` handler method that sends a 204 response.

```sql
DELETE FROM products
WHERE product_id = $1
```


## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/tests/product_test.go
Modified internal/product/product.go
Modified internal/product/product_test.go
```