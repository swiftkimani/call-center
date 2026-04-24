-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY full_name
LIMIT $1 OFFSET $2;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW()
WHERE token_hash = $1;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens SET revoked_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL;
