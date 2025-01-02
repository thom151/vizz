-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=?;


-- name: Reset :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users SET email = ?, hashed_password = ?, updated_at = ?
WHERE id = ?
RETURNING *;
