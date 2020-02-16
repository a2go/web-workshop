# 13. Error Statuses

Not every error is an internal server error.


## Tasks:

- Add a custom error type that knows about HTTP status codes.
- Make `web.Decode` return an error with status "400 Bad Request".
- Modify the middleware function to detect this case and use the provided status code.


## File Changes:

```
Modified internal/platform/web/errors.go
Modified internal/platform/web/request.go
Modified internal/platform/web/response.go
Modified internal/platform/web/web.go
```