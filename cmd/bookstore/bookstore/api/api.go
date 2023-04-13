package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"libdb.so/go-workshop/cmd/bookstore/bookstore"
)

// BookstoreHandler is the handler for the bookstore API.
type BookstoreHandler struct {
	// Embed ("inherit"-ish) chi.Router so we get http.Handler for free.
	chi.Router

	// store is the datastore for the bookstore. It is lower-case so the field
	// is private (unexported).
	store bookstore.BookStorer
}

// NewBookstoreHandler creates a new BookstoreHandler.
func NewBookstoreHandler(store bookstore.BookStorer) *BookstoreHandler {
	r := chi.NewRouter()
	h := &BookstoreHandler{
		Router: r,
		store:  store,
	}

	r.Route("/books", func(r chi.Router) {
		r.Get("/", h.getBooks)
		r.Post("/", h.addBook)
		r.Get("/{isbn}", h.getBook)
		r.Put("/{isbn}", h.putBook) // like PATCH but requires all fields
		r.Delete("/{isbn}", h.deleteBook)
	})

	return h
}

func (h *BookstoreHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.store.Books()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, books)
}

func (h *BookstoreHandler) getBook(w http.ResponseWriter, r *http.Request) {
	isbn := bookstore.ISBN(chi.URLParam(r, "isbn"))
	if err := isbn.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	book, err := h.store.Book(isbn)
	if err != nil {
		if errors.Is(err, bookstore.ErrBookNotFound) { // special 404 case
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, book)
}

func (h *BookstoreHandler) addBook(w http.ResponseWriter, r *http.Request) {
	book, err := unmarshalRequest[bookstore.Book](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.AddBook(book); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
}

func (h *BookstoreHandler) putBook(w http.ResponseWriter, r *http.Request) {
	book, err := unmarshalRequest[bookstore.Book](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.UpdateBook(book); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 so successful but no content
}

func (h *BookstoreHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := bookstore.ISBN(chi.URLParam(r, "isbn"))
	if err := isbn.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.DeleteBook(isbn); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 so successful but no content
}

// unmarshalRequest is a helper function to unmarshal a request body. It calls
// Validate() automatically on the unmarshaled object if it implements the
// Validate method. This is commonly called "duck typing".
func unmarshalRequest[T any](r *http.Request) (T, error) {
	var b T
	// Note! You MUST give JSON a pointer to the object to unmarshal into.
	// This is similar to C++'s reference semantics (T&).
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		return b, err
	}

	// Do a type-assertion: does b implement the Validate method? If so, call
	// it.
	// Go quirk: we cannot type-assert on a generic type, so we have to cast it
	// to an interface first.
	if v, ok := any(b).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return b, err
		}
	}

	return b, nil
}

func writeError(w http.ResponseWriter, code int, err error) {
	type response struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response{Error: err.Error()})
}

// writeJSON is a helper function to write a JSON response.
func writeJSON(w http.ResponseWriter, b any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}
