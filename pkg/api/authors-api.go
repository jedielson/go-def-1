package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
)

func NewAuthorsApi(r *mux.Router, repository database.AuthorsRepository) {

	r.HandleFunc("/authors", GetAuthors(repository)).Methods("GET")
}

func GetAuthors(repository database.AuthorsRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		name, limit, offset := parseUrl(r)
		authors := repository.GetAll(name, limit, offset)
		if authors == nil {
			authors = []domain.Author{}
		}

		err := json.NewEncoder(w).Encode(authors)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func parseUrl(r *http.Request) (string, int, int) {
	name := r.URL.Query().Get("name")
	offset := 0
	limit := 1000

	if i, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && i > 0 {
		offset = i
	}

	if i, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && i < 1000 && i > 0 {
		limit = i
	}

	return name, offset, limit
}
