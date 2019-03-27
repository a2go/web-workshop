# 8. Logging

- Do not use any package level variables such as the global `log`.
- Pass a `*log.Logger` to dependencies.

## Notes

Log actionable events. This is separate from Metrics or Tracing.