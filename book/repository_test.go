package book_test

import (
	"context"
	"database/sql"
	"fmt"
	"go-pbt/book"
	container "go-pbt/internal"
	"log"
	"os"
	"testing"
)

const migrationPath = "../db/migrations"

var db *sql.DB

func TestMain(m *testing.M) {
	// container起動
	container, err := container.RunMySQLContainer()
	if err != nil {
		log.Fatal(err)
	}

	// マイグレーション
	if err = container.Migrate(migrationPath); err != nil {
		container.Close()
		log.Fatal(err)
	}

	db = container.DB

	code := m.Run()

	container.Close()
	os.Exit(code)
}

func TestRepository(t *testing.T) {
	ctx := context.Background()
	br := book.NewRepository(db)

	isbn := "isbn"
	author := "author"
	title := "title"

	if err := br.AddBook(ctx, isbn, title, author); err != nil {
		t.Fatal(err)
	}

	book, err := br.FindBookByIsbn(ctx, isbn)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(book)
}
