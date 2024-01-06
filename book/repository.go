package book

import (
	"context"
	"database/sql"
	"fmt"
	"go-pbt/infrastructure"
)

type BookRepository interface {
	AddBook(ctx context.Context, isbn, title, author string, options ...addBookOptions) error
	AddCopy(ctx context.Context, isbn string) error
	BorrowCopy(ctx context.Context, isbn string) error
	ReturnCopy(ctx context.Context, isbn string) error
	FindBookByAuthor(ctx context.Context, author string) ([]infrastructure.Book, error)
	FindBookByIsbn(ctx context.Context, isbn string) (infrastructure.Book, error)
	FindBookByTitle(ctx context.Context, title string) ([]infrastructure.Book, error)
}

type bookRepository struct {
	q *infrastructure.Queries
}

func NewRepository(db *sql.DB) BookRepository {
	return &bookRepository{q: infrastructure.New(db)}
}

type AddBookOption struct {
	Owned sql.NullInt32
	Avail sql.NullInt32
}

type addBookOptions func(*AddBookOption)

func WithOwned(o *AddBookOption, owned int32) {
	o.Owned = sql.NullInt32{Int32: owned, Valid: true}
}

func WithAvail(o *AddBookOption, avail int32) {
	o.Avail = sql.NullInt32{Int32: avail, Valid: true}
}

func (br *bookRepository) AddBook(ctx context.Context, isbn, title, author string, options ...addBookOptions) error {
	var op AddBookOption
	for _, option := range options {
		option(&op)
	}

	params := infrastructure.AddBookParams{
		Isbn:      isbn,
		Title:     title,
		Author:    author,
		Owned:     op.Owned,
		Available: op.Avail,
	}

	return br.q.AddBook(ctx, params)
}

func (br *bookRepository) AddCopy(ctx context.Context, isbn string) error {
	result, err := br.q.AddCopy(ctx, isbn)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("not found isbn: %s", isbn)
	}

	return nil
}

func (br *bookRepository) BorrowCopy(ctx context.Context, isbn string) error {
	result, err := br.q.BorrowCopy(ctx, isbn)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("not found isbn: %s", isbn)
	}

	return nil
}

func (br *bookRepository) ReturnCopy(ctx context.Context, isbn string) error {
	result, err := br.q.ReturnCopy(ctx, isbn)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("not found isbn: %s", isbn)
	}

	return nil
}

func (br *bookRepository) FindBookByAuthor(ctx context.Context, author string) ([]infrastructure.Book, error) {
	return br.q.FindByAuthor(ctx, "%"+author+"%")
}

func (br *bookRepository) FindBookByIsbn(ctx context.Context, isbn string) (infrastructure.Book, error) {
	return br.q.FindByIsbn(ctx, isbn)
}

func (br *bookRepository) FindBookByTitle(ctx context.Context, title string) ([]infrastructure.Book, error) {
	return br.q.FindByTitle(ctx, "%"+title+"%")
}
