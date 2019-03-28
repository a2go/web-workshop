# 29. Authentication

- Add a key type and variable for `auth.Claims`
- In package `mid` create a type `Auth` with a method `Authenticate`.
  - This method is a middleware that ensures the request has a valid token.
  - Store the claims in the request context.
- Modify the `Handle` method in the `web` package to accept optional route specific middleware.
- Add the `Authenticate` middleware to all routes except the health check and the token generator.
