# 22. Create Users

- Add a `internal/users` package with types `User` and `NewUser`.
- Add a `Create` function to the users package.
  - It should generate a password hash using bcrypt.
- Add a `internal/plaform/auth` package with constants `RoleAdmin` and `RoleUser`.
- Add a command `useradd` to the `sales-admin` program.
