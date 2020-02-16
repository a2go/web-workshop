# 11. Web Helpers

Encoding/decoding are repetitive tasks that happen in almost every handler.

Move these behaviors to a dedicated package.

## Tasks:

- Define `Respond` and `Decode` functions to centralize request/response behavior.

## File Changes:

```
Modified cmd/sales-api/internal/handlers/product.go
Added    internal/platform/web/request.go
Added    internal/platform/web/response.go
```