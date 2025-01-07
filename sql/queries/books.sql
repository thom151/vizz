-- name: CreateBookEntry :one
INSERT INTO books (title, author, description, epub_path, created_at, updated_at)
VALUES (
    ?,
    ?,
    ?,
    ?,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
    )
RETURNING *;


-- name: GetBooks :many
SELECT * FROM books WHERE title LIKE '%' || :book || '%';

-- name: GetBook :one
SELECT * FROM books WHERE id=?;
