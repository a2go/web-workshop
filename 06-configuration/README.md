# 6. Configuration

- Create `config` struct in `cmd/sales-api` and `cmd/sales-admin`.
- Remove hardcoded HTTP & DB info.
- Copy the `internal/platform/conf` package into your project.
- Use `conf.Parse` to populate your config structs.
- Detect the `ErrHelpWanted` error and display `conf.Usage` in that case.
- Unlike the `sales-api` program, the `sales-admin` config should include a
  field of type `conf.Args` to capture command line arguments after flags are
  processed.

## Discussion

Configuration can come from many places. Some programs use environment
variables, command-line flags, config files, or configuration services.

## File Changes:

```
Modified cmd/sales-admin/main.go
Modified cmd/sales-api/main.go
Modified internal/platform/database/database.go
```
