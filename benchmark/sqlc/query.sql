-- name: Create :exec
INSERT INTO books (isbn, title, author, genre, quantity, publicized_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateMany :copyfrom
INSERT INTO books (isbn, title, author, genre, quantity, publicized_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateReturningID :one
INSERT INTO books (isbn, title, author, genre, quantity, publicized_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: Update :exec
UPDATE books
SET isbn = $1,
    title = $2,
    author = $3,
    genre = $4,
    quantity = $5,
    publicized_at = $6
WHERE id = $7;

-- name: Delete :exec
DELETE FROM books WHERE id = $1;

-- name: Get :one
SELECT * FROM books WHERE id = $1 ;

-- name: ListPaginating :many
SELECT * FROM books WHERE id > $1 LIMIT $2;
