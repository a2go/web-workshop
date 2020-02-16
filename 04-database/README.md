# 4. Database

- Start a PostgreSQL server running in a docker container
- Setup schema
- Seed database
- Connect to database from service
- Remove hardcoded product list and replace with SQL query

## Notes:

- Executing schema changes requires elevated privileges. The normal API service
  should be running as a DB user with less access.
- Using `SELECT *` has problems.

```
## Start postgres:
docker-compose up -d

## Create the schema and insert some seed data.
go build
./garagesale migrate
./garagesale seed

## Run the app then make requests.
./garagesale
```

## Links:

- [Surprises, Antipatterns and Limitations (of `database/sql`)](http://go-database-sql.org/surprises.html)


## File Changes:

```
Added    docker-compose.yaml
Modified main.go
Added    schema/migrate.go
Added    schema/seed.go
```

## Dependency Changes:

```
+ 	github.com/GuiaBolso/darwin v0.0.0-20170210191649-86919dfcf808
+ 	github.com/jmoiron/sqlx v1.2.0
+ 	github.com/lib/pq v1.1.1
```