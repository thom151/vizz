-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    DATETIME('now'),
    DATETIME('now'),
    $2,
    $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1;


-- name: Reset :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3, updated_at = DATETIME('now')
WHERE id = $1
RETURNING *;
