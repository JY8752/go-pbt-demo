-- 本を追加する
-- name: AddBook :exec
INSERT INTO books (isbn, title, author, owned, available)
VALUES (?, ?, ?, ?, ?);

-- 既存の本を1冊追加する 
-- name: AddCopy :exec
UPDATE books SET
  owned = owned + 1,
  available = available + 1 
WHERE 
  isbn = ?;

-- 本を1冊借りる
-- name: BorrowCopy :exec
UPDATE books SET available = available - 1 WHERE isbn = ? AND available > 0;

-- 本を返却する
-- name: ReturnCopy :exec
UPDATE books SET available = available + 1 WHERE isbn = ?;

-- 本を見つける
-- name: FindByAuthor :many
SELECT * FROM books WHERE author LIKE ?;

-- name: FindByIsbn :one
SELECT * FROM books WHERE isbn = ?;

-- name: FindByTitle :one
SELECT * FROM books WHERE title LIKE ?; 