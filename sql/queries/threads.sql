
-- name: CreateThread :one
INSERT INTO threads (created_at, updated_at, thread_id, user_id, book_id)
VALUES (
   CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    ?,
    ?,
    ?
    )
RETURNING *;

-- name: GetThread :one
SELECT * FROM threads WHERE user_id = ? AND book_id = ?;
