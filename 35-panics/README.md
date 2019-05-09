# 35. Panics

Applications should not use panics for control flow or error handling. It is
possible, however, that a developer might accidentally create a panic for
certain requests. The default behavior for `net/http` is to terminate the
request without responding the client. It also prints a stack trace to the
global default logger.

- Add a middleware to recover from panics and turn them into errors. This
  allows them to be tracked in the Metrics middleware and ensures the client
  will see a 500 response.


## File Changes:

```
Added    internal/mid/panics.go
```