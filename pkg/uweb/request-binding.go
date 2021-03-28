package uweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
)

type URLSegment int

const (
	Path URLSegment = iota
	Query
)

func BindGetAuthorsRequest(r *http.Request) (string, int, int) {

	name := r.URL.Query().Get("name")
	offset := FromQuery(r, "offset", 0, ValidateOffsetQuery)
	limit := FromQuery(r, "limit", 1000, ValidateLimitQuery)

	return name, offset, limit
}

func BindBookId(r *http.Request, segment URLSegment, err string) (int, error) {
	f := func(i int) bool {
		return i > 0
	}

	return FromPath(r, "id", f, errors.New(err))
}

func BindGetBooksRequest(r *http.Request) database.GetAllRequest {

	name := r.URL.Query().Get("name")
	pubYear := FromQuery(r, "publication_year", 0, func(i int) bool { return i > 1500 && i <= time.Now().Year() })
	edition := r.URL.Query().Get("edition")
	author := FromQuery(r, "author", 0, func(i int) bool { return i > 0 })
	limit := FromQuery(r, "limit", 1000, ValidateLimitQuery)
	offset := FromQuery(r, "offset", 0, ValidateOffsetQuery)
	return database.GetAllRequest{
		Name:            name,
		PublicationYear: pubYear,
		Edition:         edition,
		Author:          author,
		Limit:           limit,
		Offset:          offset,
	}
}

func BindCreateBookRequest(r *http.Request) (domain.Book, error) {
	var book domain.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		return book, errors.New("Invalid request payload")
	}

	defer r.Body.Close()
	return book, nil
}

func ValidateLimitQuery(i int) bool {
	return i < 1000 && i > 0
}

func ValidateOffsetQuery(i int) bool {
	return i > 0
}
