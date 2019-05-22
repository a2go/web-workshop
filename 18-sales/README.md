# 18. Sales

- Add a second model to represent a `Sale`.
  - This model is part of the "Products" domain so it does not need a new package.
- Add a migration for the new table and some seed data.
- Add two fields Sold and Revenue to the `Product` type.
- Modify the `List` and `Retrieve` queries to populate the Sold and Revenue fields.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/tests/product_test.go
Modified internal/product/models.go
Modified internal/product/product.go
Added    internal/product/sales.go
Added    internal/product/sales_test.go
Modified internal/schema/migrations.go
Modified internal/schema/seeds.go
```