// Package bookstore describes the models for a bookstore.
package bookstore

import (
	"fmt"
	"regexp"
)

// Cents represents a price in USD cents.
type Cents int

// String returns a string representation of the price. If a Price is used in
// fmt.Println, the String method will be called automatically.
func (p Cents) String() string {
	dollar := p / 100
	cents := p % 100
	return fmt.Sprintf("$%d.%02d", dollar, cents)
}

// ISBN represents a book's ISBN in ISBN-13 format.
type ISBN string

var isbnRe = regexp.MustCompile(`^\d{3}-\d+$`)

// Validate validates that the ISBN is valid. It returns an error if the ISBN is
// not.
func (isbn ISBN) Validate() error {
	if !isbnRe.MatchString(string(isbn)) {
		return fmt.Errorf("invalid ISBN: %s", isbn)
	}
	return nil
}

// Book represents a book.
type Book struct {
	ISBN   ISBN   `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  Cents  `json:"price"`
}

// Validate validates that the book is valid. It returns an error if the book is
// not.
func (b Book) Validate() error {
	if err := b.ISBN.Validate(); err != nil {
		return err
	}
	if b.Title == "" {
		return fmt.Errorf("title is required")
	}
	if b.Author == "" {
		return fmt.Errorf("author is required")
	}
	if b.Price < 0 {
		return fmt.Errorf("price must be positive")
	}
	return nil
}

// ErrBookNotFound is returned when a book is not found.
var ErrBookNotFound = fmt.Errorf("book not found")

// BookStorer represents a storage for books.
type BookStorer interface {
	// Book retrieves a book by its ISBN.
	Book(isbn ISBN) (Book, error)
	// Books retrieves all books.
	Books() ([]Book, error)
	// AddBook adds a new book.
	AddBook(b Book) error
	// UpdateBook updates an existing book.
	UpdateBook(b Book) error
	// DeleteBook deletes a book by its ISBN.
	DeleteBook(isbn ISBN) error
}
