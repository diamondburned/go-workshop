package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"
	"libdb.so/go-workshop/cmd/bookstore-template/bookstore"
)

// TestBookstoreHandler tests... BookstoreHandler. The convention is to have
// TestX, where X is the name of the type or function you're testing. Go would
// love you for doing that.
func TestBookstoreHandler(t *testing.T) {
	// Make a helper function that creates a new BookstoreHandler with a
	// mockBookStorer. This way we can reuse it in multiple tests.
	newServer := func(t *testing.T, mock *mockBookStorer) *httptest.Server {
		server := httptest.NewServer(NewBookstoreHandler(mock))
		t.Cleanup(server.Close) // stop the server when the test is done
		return server
	}

	// Define some books in stock for us to test with.
	books := []bookstore.Book{
		{
			ISBN:   "978-1453704424",
			Title:  "The Communist Manifesto",
			Author: "Karl Marx, Friedrich Engels",
			Price:  0, // free!
		},
		{
			ISBN:   "978-0134190440",
			Title:  "The Go Programming Language (1st Edition)",
			Author: "Alan A. A. Donovan, Brian W. Kernighan",
			Price:  3599, // $35.99
		},
	}

	// create subtests! we'll make one subtest per route that we want to
	// test on.

	t.Run("getBook", func(t *testing.T) {
		store := newMockBookStorer(books)
		server := newServer(t, store)

		r, err := server.Client().Get(server.URL + "/books/978-0134190440")
		assert.NoError(t, err, "GET error")
		assert.Equal(t, r.StatusCode, 200, "GET status code")

		gotBook := unmarshalJSON[bookstore.Book](t, r)
		assert.Equal(t, gotBook, books[1], "GET book")
	})

	t.Run("addBook", func(t *testing.T) {
		store := newMockBookStorer(books)
		server := newServer(t, store)

		addingBook := bookstore.Book{
			ISBN:   "978-1492052593",
			Title:  "Programming Rust: Fast, Safe Systems Development (2nd Edition)",
			Author: "Jim Blandy, Jason Orendorff",
			Price:  3847, // $38.47
		}

		bookJSON, err := json.Marshal(addingBook)
		assert.NoError(t, err, "marshal error")

		r, err := server.Client().Post(
			server.URL+"/books",
			"application/json", bytes.NewReader(bookJSON))
		assert.NoError(t, err, "POST error")
		assert.Equal(t, r.StatusCode, 201, "addBook POST status code")

		r2, err := server.Client().Get(server.URL + "/books/978-1492052593")
		assert.NoError(t, err, "GET error")
		assert.Equal(t, r2.StatusCode, 200, "GET status code")

		gotBook := unmarshalJSON[bookstore.Book](t, r2)
		assert.Equal(t, gotBook, addingBook, "GET book")
	})

	t.Run("getBooks", func(t *testing.T) {
		store := newMockBookStorer(books)
		server := newServer(t, store)

		r, err := server.Client().Get(server.URL + "/books")
		assert.NoError(t, err, "GET error")
		assert.Equal(t, r.StatusCode, 200, "GET status code")

		gotBooks := unmarshalJSON[[]bookstore.Book](t, r)
		assert.Equal(t, gotBooks, books, "GET books")
	})

	t.Run("putBook", func(t *testing.T) {
		store := newMockBookStorer(books)
		server := newServer(t, store)

		puttingBook := bookstore.Book{
			ISBN:   "978-0134190440",
			Title:  "The Go Programming Language (1st Edition)",
			Author: "Alan A. A. Donovan, Brian W. Kernighan",
			Price:  3000, // $30.00, discounted!
		}

		bookJSON, err := json.Marshal(puttingBook)
		assert.NoError(t, err, "marshal error")

		req, err := http.NewRequest("PUT",
			server.URL+"/books/978-0134190440", bytes.NewReader(bookJSON))
		assert.NoError(t, err, "NewRequest PUT error")

		r, err := server.Client().Do(req)
		assert.NoError(t, err, "PUT error")
		assert.Equal(t, r.StatusCode, 204, "PUT status code")

		r2, err := server.Client().Get(server.URL + "/books/978-0134190440")
		assert.NoError(t, err, "GET error")
		assert.Equal(t, r2.StatusCode, 200, "GET status code")

		gotBook := unmarshalJSON[bookstore.Book](t, r2)
		assert.Equal(t, gotBook, puttingBook, "GET book")
	})

	t.Run("deleteBook", func(t *testing.T) {
		store := newMockBookStorer(books)
		server := newServer(t, store)

		req, err := http.NewRequest("DELETE", server.URL+"/books/978-0134190440", nil)
		assert.NoError(t, err, "NewRequest DELETE error")

		r, err := server.Client().Do(req)
		assert.NoError(t, err, "DELETE error")
		assert.Equal(t, r.StatusCode, 204, "DELETE status code")

		r2, err := server.Client().Get(server.URL + "/books/978-0134190440")
		assert.NoError(t, err, "GET error")
		assert.Equal(t, r2.StatusCode, 404, "GET status code")
	})
}

func unmarshalJSON[T any](t *testing.T, r *http.Response) T {
	t.Helper()
	defer r.Body.Close()

	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	assert.NoError(t, err, "unmarshalJSON error")

	return v
}

// mockBookStorer is a fake BookStorer that we can use for testing. It holds and
// serves dummy data so we don't have to spin up a real database. It is NOT
// concurrently safe. Do NOT use it in multiple goroutines.
type mockBookStorer struct {
	books []bookstore.Book
}

func newMockBookStorer(books []bookstore.Book) *mockBookStorer {
	// Ensure all our books are valid or something. Panic if not, since it's a
	// programmer error. Panicking will cause everything to explode, which warns
	// the programmer.
	for _, book := range books {
		if err := book.Validate(); err != nil {
			panic(err)
		}
	}

	return &mockBookStorer{
		books: books,
	}
}

// Books returns all the books in the store.
func (m *mockBookStorer) Books() ([]bookstore.Book, error) {
	return append([]bookstore.Book{}, m.books...), nil
}

// bookIndex searches a book by its ISBN and returns its index in the slice.
// It is a helper function and is not exported (lowercase name). -1 is returned
// if the book is not found.
func (m *mockBookStorer) bookIndex(isbn bookstore.ISBN) (int, error) {
	for i, book := range m.books {
		if book.ISBN == isbn {
			return i, nil
		}
	}
	return -1, bookstore.ErrBookNotFound
}

// Book returns the book with the given ISBN.
func (m *mockBookStorer) Book(isbn bookstore.ISBN) (bookstore.Book, error) {
	i, err := m.bookIndex(isbn)
	if err != nil {
		return bookstore.Book{}, err
	}
	return m.books[i], nil
}

// AddBook adds a book to the store.
func (m *mockBookStorer) AddBook(book bookstore.Book) error {
	if i, _ := m.bookIndex(book.ISBN); i != -1 {
		return errors.New("book already exists")
	}
	m.books = append(m.books, book)
	return nil
}

// UpdateBook updates a book in the store.
func (m *mockBookStorer) UpdateBook(book bookstore.Book) error {
	i, err := m.bookIndex(book.ISBN)
	if err != nil {
		return err
	}

	m.books[i] = book
	return nil
}

// DeleteBook deletes a book from the store.
func (m *mockBookStorer) DeleteBook(isbn bookstore.ISBN) error {
	i, err := m.bookIndex(isbn)
	if err != nil {
		return err
	}

	//
	m.books = append(m.books[:i], m.books[i+1:]...)
	return nil
}
