# 6. Configuration

- Create `config` struct in `cmd/sales-api` and `cmd/sales-admin`.
- Remove hardcoded HTTP & DB info.
- Use `github.com/kelseyhightower/envconfig` to parse the environment.
- Define `flag.Usage` to be a function that calls `envconfig.Usage` to print expected environment variables.
- Add a flag `-config-only` that prints the config being used.


## Discussion

Configuration can come from many places. Some programs use environment
variables, command-line flags, config files, or configuration services.