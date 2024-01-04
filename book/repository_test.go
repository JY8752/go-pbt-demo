package book_test

import (
	"context"
	"database/sql"
	"fmt"
	"go-pbt/book"
	container "go-pbt/internal"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// ジェネレーター
func title() *rapid.Generator[string] {
	return rapid.String()
}

func author() *rapid.Generator[string] {
	return rapid.String()
}

func isbn() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		a := rapid.OneOf(rapid.Just("978"), rapid.Just("979")).Draw(t, "isbn-a")
		b := strconv.Itoa(rapid.IntRange(0, 9999).Draw(t, "isbn-b"))
		c := strconv.Itoa(rapid.IntRange(0, 9999).Draw(t, "isbn-c"))
		d := strconv.Itoa(rapid.IntRange(0, 999).Draw(t, "isbn-d"))
		e := rapid.StringOfN(
			rapid.RuneFrom([]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'X'}),
			1, 1, 1,
		).Draw(t, "isbn-e")

		return strings.Join([]string{a, b, c, d, e}, "-")
	})
}

func TestExample(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(isbn().Example(i))
	}
}

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
