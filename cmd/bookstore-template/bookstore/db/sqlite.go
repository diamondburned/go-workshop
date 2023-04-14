package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"libdb.so/go-workshop/cmd/bookstore-template/bookstore"
	"libdb.so/go-workshop/cmd/bookstore-template/bookstore/db/sqlite"

	// Import the SQLite driver.
	_ "modernc.org/sqlite"
)

// SQLite is a BookStorer that uses SQLite.
type SQLite struct {
	db *sql.DB
	q  *sqlite.Queries
}

// Assure that SQLite implements BookStorer.
var _ bookstore.BookStorer = (*SQLite)(nil)

// NewSQLite creates a new SQLite BookStorer.
func NewSQLite(url string) (*SQLite, error) {
	db, err := sql.Open("sqlite", url)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(sqlite.Schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &SQLite{
		db: db,
		q:  sqlite.New(db),
	}, nil
}

func (db *SQLite) Books() ([]bookstore.Book, error) {
	books, err := db.q.GetBooks(context.TODO())
	if err != nil {
		return nil, sqliteError(err)
	}
	var result []bookstore.Book
	for _, book := range books {
		result = append(result, bookstore.Book{
			ISBN:   bookstore.ISBN(book.Isbn),
			Title:  book.Title,
			Author: book.Author,
			Price:  bookstore.Cents(book.Price),
		})
	}
	return result, nil
}

func (db *SQLite) AddBook(b bookstore.Book) error {
	err := db.q.AddBook(context.TODO(), sqlite.AddBookParams{
		Isbn:   string(b.ISBN),
		Title:  b.Title,
		Author: b.Author,
		Price:  int64(b.Price),
	})
	return sqliteError(err)
}

func (db *SQLite) UpdateBook(b bookstore.Book) error {
	err := db.q.UpdateBook(context.TODO(), sqlite.UpdateBookParams{
		Isbn:   string(b.ISBN),
		Title:  b.Title,
		Author: b.Author,
		Price:  int64(b.Price),
	})
	return sqliteError(err)
}

func (db *SQLite) DeleteBook(isbn bookstore.ISBN) error {
	err := db.q.DeleteBook(context.TODO(), string(isbn))
	return sqliteError(err)
}

func sqliteError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return bookstore.ErrBookNotFound
	}
	return err
}
