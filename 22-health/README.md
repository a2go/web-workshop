# 22. Health Checks

- Add a `/v1/health` endpoint that returns 200 when the database is ready.


## File Changes:

```
Added    cmd/sales-api/internal/handlers/check.go
Modified cmd/sales-api/internal/handlers/routes.go
Modified internal/platform/database/database.go
```