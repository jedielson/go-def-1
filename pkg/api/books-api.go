package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/jedielson/bookstore/pkg/uweb"
)

func NewBooksApi(r *mux.Router, repository database.BooksRepository) {

	r.HandleFunc("/books", GetBooks(repository)).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", GetBook(repository)).Methods(http.MethodGet)

	r.HandleFunc("/books", CreateBook(repository)).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", UpdateBook(repository)).Methods(http.MethodPut)
	r.HandleFunc("/books/{id}", DeleteBook(repository)).Methods(http.MethodDelete)
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

func CreateBook(repository database.BooksRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		book, errors := uweb.BindCreateBookRequest(r)

		if errors != nil {
			uweb.ToJson(w, nil, errors)
			return
		}

		id, err := repository.Create(book)
		if err != nil {
			uweb.ToJson(w, nil, err)
		}

		uweb.ToJson(w, id)
	}
}

func UpdateBook(repository database.BooksRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := uweb.BindBookId(r, uweb.Path, IdError)
		if err != nil {
			uweb.ToJson(w, nil, err)
			return
		}

		book, errors := uweb.BindCreateBookRequest(r)

		if errors != nil {
			uweb.ToJson(w, nil, errors)
		}

		err = repository.Update(id, book)

		if err != nil {
			uweb.ToJson(w, nil, err)
		}

		uweb.ToJson(w, id)
	}
}

func DeleteBook(repository database.BooksRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uweb.BindBookId(r, uweb.Path, IdError)
		if err != nil {
			uweb.ToJson(w, nil, err)
			return
		}

		err = repository.Delete(id)

		if err != nil {
			uweb.ToJson(w, nil, err)
		}

		uweb.ToJson(w, id)

	}
}
