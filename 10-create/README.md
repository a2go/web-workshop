# 10. Create

- In the `product` package define a type `NewProduct` with fields for name, cost, and quantity.
- Add a `Create` DB function that takes a `NewProduct` and returns a `*Product`.
- Add a `Create` POST handler decodes the request body into a `NewProduct` and calls `products.Create`.

## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified internal/product/models.go
Modified internal/product/product.go
```

## Dependency Changes:

```
+ 	github.com/google/uuid v1.1.1
```