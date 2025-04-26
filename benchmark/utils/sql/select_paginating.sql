-- selectPaginating
-- $1 Cursor
-- $2 Limit
SELECT * FROM books WHERE id > $1 LIMIT $2;
