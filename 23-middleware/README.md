# 23. Middleware

- Define a type in `web` for Middleware. It is a function that takes a `web.Handler` and returns a new `web.Handler`.
- Modify the signature of `web.New` to take a variable number of middleware functions.
- Call the middleware function in a loop to wrap around a final handler.
- Extract error handling from `web.go` to a middleware function in `errors.go`.
- Pass the error middleware into `web.New`.
