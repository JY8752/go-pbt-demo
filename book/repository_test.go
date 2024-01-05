package book_test

import (
	"context"
	"database/sql"
	"fmt"
	"go-pbt/book"
	container "go-pbt/internal"
	"log"
	"os"
	"slices"
	"testing"
	"unicode"

	"pgregory.net/rapid"
)

// ジェネレーター
// func notEmpty(g *rapid.Generator[string]) *rapid.Generator[string] {
// 	return g.Filter(func(v string) bool { return len(v) != 0 })
// }

// 仕様に合わせて生成する文字列は調整　今回はASCII文字列と数字から1-100文字の範囲で生成
func title() *rapid.Generator[string] {
	return rapid.StringOfN(rapid.RuneFrom(nil, unicode.ASCII_Hex_Digit), 1, 100, -1)
}

// 仕様に合わせて生成する文字列は調整　今回はASCII文字列と数字から1-100文字の範囲で生成
func author() *rapid.Generator[string] {
	return rapid.StringOfN(rapid.RuneFrom(nil, unicode.ASCII_Hex_Digit), 1, 100, -1)
}

func isbn() *rapid.Generator[string] {
	// return rapid.Custom(func(t *rapid.T) string {
	// 	a := rapid.OneOf(rapid.Just("978"), rapid.Just("979")).Draw(t, "isbn-a")
	// 	b := strconv.Itoa(rapid.IntRange(0, 9999).Draw(t, "isbn-b"))
	// 	c := strconv.Itoa(rapid.IntRange(0, 9999).Draw(t, "isbn-c"))
	// 	d := strconv.Itoa(rapid.IntRange(0, 999).Draw(t, "isbn-d"))
	// 	e := rapid.StringOfN(
	// 		rapid.RuneFrom([]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'X'}),
	// 		1, 1, 1,
	// 	).Draw(t, "isbn-e")

	// 	return strings.Join([]string{a, b, c, d, e}, "-")
	// })
	return rapid.StringMatching("(978|979)-(([0-9]|[1-9][0-9]|[1-9]{2}[0-9]|[1-9]{3}[0-9])-){2}([0-9]|[1-9][0-9]|[1-9]{2}[0-9])-[0-9X]")
}

func TestExample(t *testing.T) {
	for i := 0; i < 20; i++ {
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

func TestProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		type state struct {
			isbn   string
			author string
			title  string
		}

		ctx := context.Background()
		br := book.NewRepository(db)
		states := make([]state, 0, 100)
		isbns := make([]string, 0, 100)

		t.Repeat(map[string]func(*rapid.T){
			"AddBook": func(t *rapid.T) {
				isbn := isbn().Filter(func(v string) bool {
					return !slices.Contains(isbns, v)
				}).Draw(t, "isbn")
				author := author().Draw(t, "author")
				title := title().Draw(t, "title")
				if err := br.AddBook(ctx, isbn, title, author); err != nil {
					t.Fatal(err)
				}
				states = append(states, state{
					isbn:   isbn,
					author: author,
					title:  title,
				})
				isbns = append(isbns, isbn)
			},
			"AddCopy": func(t *rapid.T) {
				if len(states) == 0 {
					t.Skip("no books")
				}

				isbn := states[len(states)-1].isbn
				if err := br.AddCopy(ctx, isbn); err != nil {
					t.Fatalf("failed to AddCopy isbn: %s err: %s", isbn, err.Error())
				}
			},
			"BorrowCopy": func(t *rapid.T) {
				if len(states) == 0 {
					t.Skip("no books")
				}

				isbn := states[len(states)-1].isbn
				if err := br.BorrowCopy(ctx, isbn); err != nil {
					t.Fatalf("failed to BorrowCopy isbn: %s err: %s", isbn, err.Error())
				}
			},
			"ReturnCopy": func(t *rapid.T) {
				if len(states) == 0 {
					t.Skip("no books")
				}

				isbn := states[len(states)-1].isbn
				if err := br.ReturnCopy(ctx, isbn); err != nil {
					t.Fatalf("failed to ReturnCopy isbn: %s err: %s", isbn, err.Error())
				}
			},
			"FindBookByAuthor": func(t *rapid.T) {
				if len(states) == 0 {
					t.Skip("no books")
				}

				state := states[len(states)-1]
				_, err := br.FindBookByAuthor(ctx, state.author)
				if err != nil {
					t.Fatalf("failed to FindBookByAuthor isbn: %s err: %s", state.isbn, err.Error())
				}
			},
			"FindBookByTitle": func(t *rapid.T) {
				if len(states) == 0 {
					t.Skip("no books")
				}

				state := states[len(states)-1]
				_, err := br.FindBookByTitle(ctx, state.title)
				if err != nil {
					t.Fatalf("failed to FindBookByTitle isbn: %s err: %s", state.isbn, err.Error())
				}
			},
			"FindBookByIsbn": func(t *rapid.T) {
				if len(states) == 0 {
					t.Skip("no books")
				}

				isbn := states[len(states)-1].isbn
				_, err := br.FindBookByIsbn(ctx, isbn)
				if err != nil {
					t.Fatalf("failed to FindBookByIsbn isbn: %s err: %s", isbn, err.Error())
				}
			},
		})
	})
}
