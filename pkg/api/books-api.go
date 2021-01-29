package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
)

func NewBooksApi(r *mux.Router, repository database.BooksRepository) {

	r.HandleFunc("/books", GetBooks(repository)).Methods("GET")
}

func GetBooks(repository database.BooksRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// name, limit, offset := parseUrl(r)
		books := repository.GetAll(database.GetAllRequest{})
		if books == nil {
			books = []domain.Book{}
		}

		err := json.NewEncoder(w).Encode(books)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusBadGateway)
	}
}
