-- name: GetCustomerByID :one
SELECT * FROM customers WHERE id = $1;

-- name: GetCustomerByPhone :one
SELECT * FROM customers WHERE phone_number = $1;

-- name: CreateCustomer :one
INSERT INTO customers (phone_number, full_name, email, tags)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET full_name = $2, email = $3, tags = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkCustomerDNC :exec
UPDATE customers SET dnc_listed = true, updated_at = NOW() WHERE id = $1;

-- name: SearchCustomers :many
SELECT * FROM customers
WHERE
  full_name ILIKE '%' || $1 || '%'
  OR phone_number ILIKE '%' || $1 || '%'
  OR email ILIKE '%' || $1 || '%'
ORDER BY full_name
LIMIT $2 OFFSET $3;
