# 18. Sales

- Add a second model to represent a `Sale`.
  - This model is part of the "Products" domain so it does not need a new package.
- Add a migration for the new table and some seed data.
- Add two fields Sold and Revenue to the `Product` type.
- Modify the `List` and `Get` queries to populate the Sold and Revenue fields.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/products.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/tests/products_test.go
Modified internal/platform/database/schema/migrations.go
Modified internal/platform/database/schema/seeds.go
Modified internal/products/products.go
Added    internal/products/sales.go
Added    internal/products/sales_test.go
```