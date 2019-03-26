# 25. Login

- Add a `keygen` command to `sales-admin` which can generate a x509 private key.
- Add a `Claims` type to the `auth` package.
- Add an `Authenticator` type to the `auth` package which generates a JWT for `Claims`.
- Add an `Authenticate` function to the `users` package which finds a user by email and verifies their password.
- Add a handler that identifies a user from Basic Auth and responds with a Token.
