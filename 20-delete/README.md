# 20. Delete

- Add `Delete` function in `internal/products`,
- Add route for `DELETE` `/v1/products/{id}`.
- Add `Delete` handler method that sends a 204 response.

```sql
DELETE FROM products
WHERE product_id = $1
```
