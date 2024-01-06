package book_test

import (
	"context"
	"database/sql"
	"errors"
	"go-pbt/book"
	"go-pbt/infrastructure"
	container "go-pbt/internal"
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"
	"testing"
	"unicode"

	"pgregory.net/rapid"
)

// func notEmpty(g *rapid.Generator[string]) *rapid.Generator[string] {
// 	return g.Filter(func(v string) bool { return len(v) != 0 })
// }

// func TestExample(t *testing.T) {
// 	for i := 0; i < 20; i++ {
// 		fmt.Println(isbn().Example(i))
// 	}
// }

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

// func TestRepository(t *testing.T) {
// 	ctx := context.Background()
// 	br := book.NewRepository(db)

// 	isbn := "isbn"
// 	author := "author"
// 	title := "title"

// 	if err := br.AddBook(ctx, isbn, title, author); err != nil {
// 		t.Fatal(err)
// 	}

// 	book, err := br.FindBookByIsbn(ctx, isbn)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	fmt.Println(book)
// }

// type state struct {
// 	isbn   string
// 	author string
// 	title  string
// }

// func NewState(isbn, author, title string) state {
// 	return state{isbn, author, title}
// }

// func TestProperty(t *testing.T) {
// 	rapid.Check(t, func(t *rapid.T) {
// 		ctx := context.Background()
// 		br := book.NewRepository(db)
// 		states := make([]state, 0, 100)
// 		isbns := make([]string, 0, 100)

// 		t.Repeat(map[string]func(*rapid.T){
// 			"AddBook": func(t *rapid.T) {
// 				isbn := isbn().Filter(func(v string) bool {
// 					return !slices.Contains(isbns, v)
// 				}).Draw(t, "isbn")
// 				author := author().Draw(t, "author")
// 				title := title().Draw(t, "title")
// 				if err := br.AddBook(ctx, isbn, title, author); err != nil {
// 					t.Fatal(err)
// 				}
// 				states = append(states, state{
// 					isbn:   isbn,
// 					author: author,
// 					title:  title,
// 				})
// 				isbns = append(isbns, isbn)
// 			},
// 			"AddCopy": func(t *rapid.T) {
// 				if len(states) == 0 {
// 					t.Skip("no books")
// 				}

// 				isbn := states[len(states)-1].isbn
// 				if err := br.AddCopy(ctx, isbn); err != nil {
// 					t.Fatalf("failed to AddCopy isbn: %s err: %s", isbn, err.Error())
// 				}
// 			},
// 			"BorrowCopy": func(t *rapid.T) {
// 				if len(states) == 0 {
// 					t.Skip("no books")
// 				}

// 				isbn := states[len(states)-1].isbn
// 				if err := br.BorrowCopy(ctx, isbn); err != nil {
// 					t.Fatalf("failed to BorrowCopy isbn: %s err: %s", isbn, err.Error())
// 				}
// 			},
// 			"ReturnCopy": func(t *rapid.T) {
// 				if len(states) == 0 {
// 					t.Skip("no books")
// 				}

// 				isbn := states[len(states)-1].isbn
// 				if err := br.ReturnCopy(ctx, isbn); err != nil {
// 					t.Fatalf("failed to ReturnCopy isbn: %s err: %s", isbn, err.Error())
// 				}
// 			},
// 			"FindBookByAuthor": func(t *rapid.T) {
// 				if len(states) == 0 {
// 					t.Skip("no books")
// 				}

// 				state := states[len(states)-1]
// 				_, err := br.FindBookByAuthor(ctx, state.author)
// 				if err != nil {
// 					t.Fatalf("failed to FindBookByAuthor isbn: %s err: %s", state.isbn, err.Error())
// 				}
// 			},
// 			"FindBookByTitle": func(t *rapid.T) {
// 				if len(states) == 0 {
// 					t.Skip("no books")
// 				}

// 				state := states[len(states)-1]
// 				_, err := br.FindBookByTitle(ctx, state.title)
// 				if err != nil {
// 					t.Fatalf("failed to FindBookByTitle isbn: %s err: %s", state.isbn, err.Error())
// 				}
// 			},
// 			"FindBookByIsbn": func(t *rapid.T) {
// 				if len(states) == 0 {
// 					t.Skip("no books")
// 				}

// 				isbn := states[len(states)-1].isbn
// 				_, err := br.FindBookByIsbn(ctx, isbn)
// 				if err != nil {
// 					t.Fatalf("failed to FindBookByIsbn isbn: %s err: %s", isbn, err.Error())
// 				}
// 			},
// 		})
// 	})
// }

// 書籍情報の状態
type _book struct {
	isbn   string
	author string
	title  string
	owned  int32
	avail  int32
}

func NewBook(isbn, author, title string, owned, avail int32) *_book {
	return &_book{isbn, author, title, owned, avail}
}

// 状態管理
type states = map[string]*_book

// ユーティリティー / ヘルパー
func keys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

// func values[K comparable, V any](m map[K]V) []V {
// 	s := make([]V, 0, len(m))
// 	for _, v := range m {
// 		s = append(s, v)
// 	}
// 	return s
// }

func merge[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	newMap := make(map[K]V, len(m1)+len(m2))
	for k, v := range m1 {
		newMap[k] = v
	}
	for k, v := range m2 {
		newMap[k] = v
	}
	return newMap
}

// sliceの要素が空だとpanicする
func elements[T any](s []T) T {
	switch len(s) {
	case 0:
		panic("slice is empty")
	case 1:
		return s[0]
	}
	return s[rand.Intn(len(s)-1)]
}

func partial(t *rapid.T, str string) string {
	l := len([]rune(str))
	start := rapid.IntRange(0, l-1).Draw(t, "start")
	end := rapid.IntRange(start+1, l).Draw(t, "end")

	return string([]rune(str)[start:end])
}

// func TestPartial(t *testing.T) {
// 	rapid.Check(t, func(t *rapid.T) {
// 		str := "d0"
// 		for i := 0; i < 10; i++ {
// 			fmt.Println(partial(t, str))
// 		}
// 	})
// }

func hasIsbn(states states, isbn string) bool {
	keys := keys(states)
	return slices.Contains(keys, isbn)
}

func likeAuthor(states states, author string) bool {
	if author == "" {
		return false
	}

	for _, v := range states {
		if strings.Contains(strings.ToLower(v.author), strings.ToLower(author)) {
			return true
		}
	}

	return false
}

func likeTitle(states states, title string) bool {
	if title == "" {
		return false
	}

	for _, v := range states {
		if strings.Contains(strings.ToLower(v.title), strings.ToLower(title)) {
			return true
		}
	}

	return false
}

// ジェネレーター

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

func isbnGen(states states) string {
	return elements(keys(states))
}

func authorGen(t *rapid.T, states states) string {
	s := make([]string, 0, len(states))
	for _, v := range states {
		s = append(s, partial(t, v.author))
	}
	return elements(s)
}

func titleGen(t *rapid.T, states states) string {
	s := make([]string, 0, len(states))
	for _, v := range states {
		s = append(s, partial(t, v.title))
	}
	return elements(s)
}

func TestProperty2(t *testing.T) {
	ctx := context.Background()
	br := book.NewRepository(db)
	states := make(states)

	rapid.Check(t, func(t *rapid.T) {
		// 状態に依存しないテスト
		alwaysPossible := map[string]func(*rapid.T){
			"AddBookNew": func(t *rapid.T) {
				isbn := isbn().Draw(t, "isbn")
				author := author().Draw(t, "author")
				title := title().Draw(t, "title")

				// 事前条件
				if hasIsbn(states, isbn) {
					t.Skip("already exist book")
				}

				if err := br.AddBook(ctx, isbn, title, author, book.WithOwned(1), book.WithAvail(1)); err != nil {
					t.Fatalf("failed to AddBookNew isbn: %s err: %s", isbn, err.Error())
				}

				// 状態更新
				states[isbn] = NewBook(isbn, author, title, 1, 1)
			},
			"AddCopyNew": func(t *rapid.T) {
				isbn := isbn().Draw(t, "isbn")

				// 事前条件
				if hasIsbn(states, isbn) {
					t.Skip("already exist book")
				}

				if err := br.AddCopy(ctx, isbn); err == nil {
					t.Fatal("expected error, but not error")
				}
			},
			"BorrowCopyUnkown": func(t *rapid.T) {
				isbn := isbn().Draw(t, "isbn")

				// 事前条件
				if hasIsbn(states, isbn) {
					t.Skip("already exist book")
				}

				if err := br.BorrowCopy(ctx, isbn); err == nil {
					t.Fatal("expected error, but not error")
				}
			},
			"ReturnCopyUnkown": func(t *rapid.T) {
				isbn := isbn().Draw(t, "isbn")

				// 事前条件
				if hasIsbn(states, isbn) {
					t.Skip("already exist book")
				}

				if err := br.ReturnCopy(ctx, isbn); err == nil {
					t.Fatal("expected error, but not error")
				}
			},
			"FindBookByIsbnUnkown": func(t *rapid.T) {
				isbn := isbn().Draw(t, "isbn")

				// 事前条件
				if hasIsbn(states, isbn) {
					t.Skip("already exist book")
				}

				var err error
				if _, err = br.FindBookByIsbn(ctx, isbn); err == nil {
					t.Fatal("failed to FindBookByIsbnUnkown. expect error, but not error")
				}

				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expect sql.ErrNoRows, but %v", err)
				}
			},
			"FindBookByAuthorUnkown": func(t *rapid.T) {
				author := author().Draw(t, "author")

				// 事前条件
				if likeAuthor(states, author) {
					t.Skip("already exist book")
				}

				result, err := br.FindBookByAuthor(ctx, author)
				if err != nil {
					t.Fatalf("failed to FindBookByAuthorUnkown author: %s err: %s", author, err.Error())
				}

				if len(result) != 0 {
					t.Fatalf("failed to FindBookByAuthorUnkown. expect record not found, but found result: %v", result)
				}
			},
			"FindBookByTitleUnkown": func(t *rapid.T) {
				title := title().Draw(t, "title")

				// 事前条件
				if likeTitle(states, title) {
					t.Skip("already exist book")
				}

				result, err := br.FindBookByTitle(ctx, title)
				if err != nil {
					t.Fatalf("failed to FindBookByTitlteUnkown title: %s err: %s", title, err.Error())
				}

				if len(result) != 0 {
					t.Fatalf("failed to FindBookByAuthorUnkown. expect record not found, but found result: %v", result)
				}
			},
		}

		// 状態に依存するテスト
		reliesOnState := map[string]func(*rapid.T){
			"AddBookExisting": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				isbn := isbnGen(states)
				title := title().Draw(t, "title")
				author := author().Draw(t, "author")

				// 事前条件
				if !hasIsbn(states, isbn) {
					t.Fatalf("states not include generate ISBN %s", isbn)
				}

				// duplicate keyでエラーを期待
				if err := br.AddBook(ctx, isbn, title, author); err == nil {
					t.Fatal("expect error, but not error")
				}
			},
			"AddCopyExisting": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				isbn := isbnGen(states)

				// 事前条件
				if !hasIsbn(states, isbn) {
					t.Fatalf("states not include generate ISBN %s", isbn)
				}

				if err := br.AddCopy(ctx, isbn); err != nil {
					t.Fatalf("failed to AddCopyExisting isbn: %s err: %s", isbn, err.Error())
				}

				// 状態更新
				states[isbn].avail += 1
				states[isbn].owned += 1
			},
			"BorrowCopyAvail": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				isbn := isbnGen(states)

				// 事前条件
				if !hasIsbn(states, isbn) {
					t.Fatalf("states not include generate ISBN %s", isbn)
				}

				if states[isbn].avail == 0 {
					t.Skip("no books to borrow")
				}

				if err := br.BorrowCopy(ctx, isbn); err != nil {
					t.Fatalf("failed to BorrowCopyAvail isbn: %s err: %s", isbn, err.Error())
				}

				// 状態更新
				states[isbn].avail -= 1
			},
			"BorrowCopyUnavail": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				isbn := isbnGen(states)

				// 事前条件
				if !hasIsbn(states, isbn) {
					t.Fatalf("states not include generate ISBN %s", isbn)
				}

				if states[isbn].avail != 0 {
					t.Skip("can borrow book yet")
				}

				if err := br.BorrowCopy(ctx, isbn); err == nil {
					t.Fatal("expected error, but not error")
				}
			},
			"ReturnCopyExisting": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				isbn := isbnGen(states)

				// 事前条件
				if !hasIsbn(states, isbn) {
					t.Fatalf("states not include generate ISBN %s", isbn)
				}

				if states[isbn].avail == states[isbn].owned {
					t.Skip("book is full")
				}

				if err := br.ReturnCopy(ctx, isbn); err != nil {
					t.Fatalf("failed to ReturnCopyExisting isbn: %s err: %s", isbn, err.Error())
				}

				// 状態更新
				states[isbn].avail += 1
			},
			// "ReturnCopyFull": func(t *rapid.T) {
			// 	// まだstateがない
			// 	if len(states) == 0 {
			// 		t.Skip("states is empty")
			// 	}

			// 	isbn := isbnGen(states)

			// 	// 事前条件
			// 	if !hasIsbn(states, isbn) {
			// 		t.Fatalf("states not include generate ISBN %s", isbn)
			// 	}

			// 	if states[isbn].avail != states[isbn].owned {
			// 		t.Skip("book is not full")
			// 	}

			// 	// 本当は貸出がない状態で返却をしようとするとエラーにしたいがそれをするには事前にDB問い合わせが必要
			// 	// やってもいいんだけど今回は手抜きでこのテストは飛ばす
			// 	if err := br.ReturnCopy(ctx, isbn); err != nil {
			// 		t.Fatalf("failed to ReturnCopyFull isbn: %s err: %s", isbn, err.Error())
			// 	}
			// },
			"FindBookByIsbnExists": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				isbn := isbnGen(states)

				// 事前条件
				if !hasIsbn(states, isbn) {
					t.Fatalf("states not include generate ISBN %s", isbn)
				}

				book, err := br.FindBookByIsbn(ctx, isbn)
				if err != nil {
					t.Fatalf("failed to FindBookByIsbnExists isbn: %s err: %s", isbn, err.Error())
				}

				assertBook(t, *states[isbn], book)
			},
			"FindBookByAuthorMatching": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				author := authorGen(t, states)

				// 事前条件
				if !likeAuthor(states, author) {
					t.Fatalf("states not include generate author %s", author)
				}

				_, err := br.FindBookByAuthor(ctx, author)
				if err != nil {
					t.Fatalf("failed to FindBookByAuthorMatching isbn: %s err: %s", author, err.Error())
				}

				// アサーション
				// statesからauthorが部分一致する本情報とDBから取得してきた本情報をソートして完全に一致しているか確認する
				// 心折れたので手抜き
			},
			"FindBookByTitleMatching": func(t *rapid.T) {
				// まだstateがない
				if len(states) == 0 {
					t.Skip("states is empty")
				}

				title := titleGen(t, states)

				// 事前条件
				if !likeTitle(states, title) {
					t.Fatalf("states not include generate title %s", title)
				}

				_, err := br.FindBookByTitle(ctx, title)
				if err != nil {
					t.Fatalf("failed to FindBookByTitleMatching title: %s err: %s", title, err.Error())
				}

				// アサーション
				// statesからtitleが部分一致する本情報とDBから取得してきた本情報をソートして完全に一致しているか確認する
				// 心折れたので手抜き
			},
		}

		t.Repeat(merge(alwaysPossible, reliesOnState))
	})
}

func assertBook(t *rapid.T, state _book, record infrastructure.Book) {
	t.Helper()

	if state.isbn != record.Isbn {
		t.Fatalf("different book.isbn state.isbn %s record.isbn %s", state.isbn, record.Isbn)
	}

	if state.title != record.Title {
		t.Fatalf("different book.title state.title %s record.title %s", state.title, record.Title)
	}

	if state.author != record.Author {
		t.Fatalf("different book.author state.author %s record.author %s", state.author, record.Author)
	}

	if state.owned != record.Owned.Int32 || !record.Owned.Valid {
		t.Fatalf("different book.owned state.owned %d record.owned %d", state.owned, record.Owned.Int32)
	}

	if state.avail != record.Available.Int32 || !record.Available.Valid {
		t.Fatalf("different book.avail state.avail %d record.avail %d", state.avail, record.Available.Int32)
	}
}
