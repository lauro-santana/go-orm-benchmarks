-- insertBooks
-- insertBook
-- $1 ISBN
-- $2 Title
-- $3 Author
-- $4 Genre
-- $5 Quantity
-- $6 Publishing date
INSERT INTO books (isbn, title, author, genre, quantity, publicized_at)
VALUES %s;
