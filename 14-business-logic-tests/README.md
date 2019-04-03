# 14. Business-Logic Tests

- Add `internal/products` test incorporating Create, Get, & List.
- Add test helper function for setting up the database.


## File Changes:

```
Added    internal/platform/database/databasetest/docker.go
Added    internal/platform/database/databasetest/setup.go
Added    internal/products/products_test.go
```

## Dependency Changes:

```
+ 	github.com/google/go-cmp v0.2.0
```