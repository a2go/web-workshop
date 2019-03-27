# 12. Erorrs

Handling application errors in a consistent and reliable way is very
repetitive. Provide some support for that in the web package.


## Tasks:

- Add a custom error type that knows about HTTP status codes.
- Define a custom signature for all handler functions that includes returning errors.
- Add a middleware function that will be ran for all handlers which deals with the returned errors.
