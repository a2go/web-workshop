# 14. Sales

- Add a second model to represent a `Sale`.
  - This model is part of the "Products" domain so it does not need a new package.
- Add a migration for the new table and some seed data.


```
# Replace cmd and internal folders of garagesale with 14-sales

docker-compose rm -f db
docker-compose up -d

go run ./cmd/sales-admin migrate
go run ./cmd/sales-admin seed     
go run ./cmd/sales-api
```