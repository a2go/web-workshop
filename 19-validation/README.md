# 19. Validation

- Copy `errors.go`, `request.go`, and `web.go` from the `web` package
    - Add the `go-playground/validator` library and dependencies.
    - Expand the error type in `web` to include specific fields.
    - Integrate validation with the `web.Decode` function.
- Add validation tags to the `NewProduct` type.


## File Changes:

```
Modified cmd/sales-api/tests/product_test.go
Modified internal/platform/web/errors.go
Modified internal/platform/web/request.go
Added    internal/platform/web/request_test.go
Modified internal/platform/web/response.go
Modified internal/platform/web/web.go
Modified internal/product/models.go
```

## Dependency Changes:

```
+ 	github.com/go-playground/locales v0.12.1
+ 	github.com/go-playground/universal-translator v0.16.0
+ 	gopkg.in/go-playground/validator.v9 v9.27.0
```