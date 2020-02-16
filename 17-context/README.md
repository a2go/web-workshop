# 17. Context

Long running operations should be given a deadline. The idiomatic way to handle cancellation is passing `context.Context` to functions that know to check for cancellation and terminate early.

- Add `context.Context` argument to `products.List`, `products.Retrieve`, and `products.Create`.
- Pass the `ctx` variable to `db.SelectContext`, `db.GetContext`, and `db.ExecContext`
- Pass the value of `r.Context()` from handlers into service functions.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Modified internal/product/product.go
Modified internal/product/product_test.go
```