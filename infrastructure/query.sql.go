// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package infrastructure

import (
	"context"
	"database/sql"
)

const addBook = `-- name: AddBook :exec
INSERT INTO books (isbn, title, author, owned, available)
VALUES (?, ?, ?, ?, ?)
`

type AddBookParams struct {
	Isbn      string
	Title     string
	Author    string
	Owned     sql.NullInt32
	Available sql.NullInt32
}

// 本を追加する
func (q *Queries) AddBook(ctx context.Context, arg AddBookParams) error {
	_, err := q.db.ExecContext(ctx, addBook,
		arg.Isbn,
		arg.Title,
		arg.Author,
		arg.Owned,
		arg.Available,
	)
	return err
}

const addCopy = `-- name: AddCopy :exec
UPDATE books SET
  owned = owned + 1,
  available = available + 1 
WHERE 
  isbn = ?
`

// 既存の本を1冊追加する
func (q *Queries) AddCopy(ctx context.Context, isbn string) error {
	_, err := q.db.ExecContext(ctx, addCopy, isbn)
	return err
}

const borrowCopy = `-- name: BorrowCopy :exec
UPDATE books SET available = available - 1 WHERE isbn = ? AND available > 0
`

// 本を1冊借りる
func (q *Queries) BorrowCopy(ctx context.Context, isbn string) error {
	_, err := q.db.ExecContext(ctx, borrowCopy, isbn)
	return err
}

const findByAuthor = `-- name: FindByAuthor :many
SELECT isbn, title, author, owned, available FROM books WHERE author LIKE ?
`

// 本を見つける
func (q *Queries) FindByAuthor(ctx context.Context, author string) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, findByAuthor, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.Isbn,
			&i.Title,
			&i.Author,
			&i.Owned,
			&i.Available,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findByIsbn = `-- name: FindByIsbn :one
SELECT isbn, title, author, owned, available FROM books WHERE isbn = ?
`

func (q *Queries) FindByIsbn(ctx context.Context, isbn string) (Book, error) {
	row := q.db.QueryRowContext(ctx, findByIsbn, isbn)
	var i Book
	err := row.Scan(
		&i.Isbn,
		&i.Title,
		&i.Author,
		&i.Owned,
		&i.Available,
	)
	return i, err
}

const findByTitle = `-- name: FindByTitle :one
SELECT isbn, title, author, owned, available FROM books WHERE title LIKE ?
`

func (q *Queries) FindByTitle(ctx context.Context, title string) (Book, error) {
	row := q.db.QueryRowContext(ctx, findByTitle, title)
	var i Book
	err := row.Scan(
		&i.Isbn,
		&i.Title,
		&i.Author,
		&i.Owned,
		&i.Available,
	)
	return i, err
}

const returnCopy = `-- name: ReturnCopy :exec
UPDATE books SET available = available + 1 WHERE isbn = ?
`

// 本を返却する
func (q *Queries) ReturnCopy(ctx context.Context, isbn string) error {
	_, err := q.db.ExecContext(ctx, returnCopy, isbn)
	return err
}