# 27. Context Part 2
## The Return of Context

We are now using middleware chains and a custom handler. We are also making
greater use of context for passing values. As we continue to expand this
program we will find relying on `request.Context` and `request.WithContext` to
be annoying.

- Add `ctx context.Context` as the first argument in the `web.Handler` type.
- Make the adapter function in `web.go`pass down the first context which was
  derived from `r.Context()`.
- Anything using `r.Context()` should instead use the passed `ctx`.
- Make everything compile.


## File Changes:

```
Modified cmd/sales-api/internal/handlers/check.go
Modified cmd/sales-api/internal/handlers/products.go
Modified internal/mid/errors.go
Modified internal/mid/logger.go
Modified internal/mid/metrics.go
Modified internal/platform/web/web.go
```