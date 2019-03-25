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
# 4. Start postgres:
docker-compose up -d

# 4. Create the schema and insert some seed data.
go build
./garagesale migrate
./garagesale seed

# 4. Run the app then make requests.
./garagesale
```
