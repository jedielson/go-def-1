package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/jedielson/bookstore/pkg/uweb"
)

func NewBooksApi(r *mux.Router, repository database.BooksRepository) {

	r.HandleFunc("/books", GetBooks(repository)).Methods("GET")
	r.HandleFunc("/books/{id}", GetBook(repository)).Methods("GET")
}

func GetBooks(repository database.BooksRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		getAllRequest := uweb.BindGetBooksRequest(r)

		books := repository.GetAll(getAllRequest)
		if books == nil {
			books = []domain.Book{}
		}

		uweb.ToJson(w, books)
	}
}

const IdError = "id is invalid"

func GetBook(repository database.BooksRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := uweb.BindBookId(r, uweb.Path, IdError)
		if err != nil {
			uweb.ToJson(w, nil, err)
			return
		}

		book, err := repository.GetBook(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		uweb.ToJson(w, book)
	}
}
