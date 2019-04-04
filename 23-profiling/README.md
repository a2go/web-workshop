# 23. Profiling

- Import `net/http/pprof` to register the pprof profiler.
  - https://golang.org/pkg/net/http/pprof/
- Launch a second HTTP service listening on a different port.
  - https://mmcloughlin.com/posts/your-pprof-is-showing

http://localhost:6060/debug/pprof/


## File Changes:

```
Modified cmd/sales-api/main.go
```