// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package infrastructure

import (
	"database/sql"
)

type Book struct {
	Isbn      string
	Title     string
	Author    string
	Owned     sql.NullInt32
	Available sql.NullInt32
}
