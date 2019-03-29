# 15. API Tests

- Add `cmd/sales-api/tests` package.
- Add tests that create a database in Docker and migrate it.
- Construct a `web.App` and run the `ServeHTTP` method with different requests.


## File Changes:

```
Added    cmd/sales-api/tests/products_test.go
```

## Dependency Changes:

```
+ 	github.com/google/go-cmp v0.2.0
```