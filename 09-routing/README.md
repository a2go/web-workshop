# 9. Routing

- Add a second endpoint: `Get`.
- Add router in `routes.go` to tie the endpoints together.
- Add `internal/platform/web` with type `App` to hold the router.
- Update `Handler` in `main`.

```
http://localhost:8000/v1/products
http://localhost:8000/v1/products/a2b0639f-2cc6-44b8-b97b-15d69dbb511e
```

## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Added    cmd/sales-api/internal/handlers/routes.go
Modified cmd/sales-api/main.go
Added    internal/platform/web/web.go
Modified internal/product/product.go
```

## Dependency Changes:

```
+ 	github.com/go-chi/chi v4.0.2+incompatible
```