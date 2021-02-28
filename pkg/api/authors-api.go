package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/jedielson/bookstore/pkg/uweb"
)

func NewAuthorsApi(r *mux.Router, repository database.AuthorsRepository) {

	r.HandleFunc("/authors", GetAuthors(repository)).Methods("GET")
}

func GetAuthors(repository database.AuthorsRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		name, limit, offset := uweb.BindGetAuthorsRequest(r)
		authors := repository.GetAll(name, limit, offset)
		if authors == nil {
			authors = []domain.Author{}
		}

		uweb.ToJson(w, authors)
	}
}
