-- updateBook
-- $1 ISBN
-- $2 Title
-- $3 Author
-- $4 Genre
-- $5 Quantity
-- $6 Publishing date
-- $7 ID
UPDATE books
SET isbn = $1,
    title = $2,
    author = $3,
    genre = $4,
    quantity = $5,
    publicized_at = $6
WHERE id = $7;
