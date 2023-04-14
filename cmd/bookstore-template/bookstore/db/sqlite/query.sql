-- TODO: implement GetBook

-- name: GetBooks :many
SELECT * FROM books;

-- name: AddBook :exec
INSERT INTO books (isbn, title, author, price) VALUES (?, ?, ?, ?);

-- name: UpdateBook :exec
UPDATE books SET title = ?, author = ?, price = ? WHERE isbn = ?;

-- name: DeleteBook :exec
DELETE FROM books WHERE isbn = ?;
