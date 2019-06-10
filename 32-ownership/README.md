# 32. Ownership

- Add the `user_id` column to the `products` table.
- When creating a product, set the User ID.
- When fetching or listing products include the User ID.
- When updating or deleting products ensure that the client either has
  `RoleAdmin` or is the owner of the specified product.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/tests/product_test.go
Modified internal/platform/tests/tests.go
Modified internal/product/models.go
Modified internal/product/product.go
Modified internal/product/product_test.go
Modified internal/product/sales_test.go
Modified internal/schema/migrate.go
```