package uweb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/stretchr/testify/suite"
)

type RequestBindingHandlerSuite struct {
	suite.Suite

	ctx context.Context
	req *http.Request
	res *httptest.ResponseRecorder
}

func (s *RequestBindingHandlerSuite) SetupTest() {
	s.ctx = context.Background()

	s.res = httptest.NewRecorder()
}

// se o nome não for informado, não deve ser usado como filtro
// se o nome for informado, deve ser usado como filtro
var testsNome = []struct {
	nome     string
	expected string
}{
	{
		nome:     "",
		expected: "",
	},
	{
		nome:     "myname",
		expected: "myname",
	},
}

func (s *RequestBindingHandlerSuite) TestNameQuery() {
	for _, n := range testsNome {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/authors?name=%s", n.nome), nil)
		nome, _, _ := BindGetAuthorsRequest(s.req)
		s.Assert().Equal(n.expected, nome)
	}
}

// se o limit não for informado, deve ser usado 1000
// se o limit não for numero deve ser usado 1000
// se o limit for maior que 1000, deve ser usado 1000
// se o limit for < 0 deve ser usado 1000
// se o limit for > 0 e < 1000 deve ser usado limit
var testsLimit = []struct {
	limit    string
	expected int
}{
	{
		limit:    "",
		expected: 1000,
	},
	{
		limit:    "limit=notANumber&&",
		expected: 1000,
	},
	{
		limit:    "limit=5000",
		expected: 1000,
	},
	{
		limit:    "limit=-1",
		expected: 1000,
	},
	{
		limit:    "limit=10",
		expected: 10,
	},
}

func (s *RequestBindingHandlerSuite) TestLimitQuery() {
	for _, n := range testsLimit {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/authors?%s", n.limit), nil)
		_, _, limit := BindGetAuthorsRequest(s.req)
		s.Assert().Equal(n.expected, limit)
	}
}

// se o offset não for informado, deve ser usado 0
// se o offset não for numero, deve ser usado 0
// se o offset for < 0 deve ser usado 0
// se o offset for > 0 deve ser usado offset
var testsOffset = []struct {
	offset   string
	expected int
}{
	{
		offset:   "",
		expected: 0,
	},
	{
		offset:   "offset=notANumber&&",
		expected: 0,
	},
	{
		offset:   "offset=-1",
		expected: 0,
	},
	{
		offset:   "offset=10",
		expected: 10,
	},
}

func (s *RequestBindingHandlerSuite) TestOffsetQuery() {
	for _, n := range testsOffset {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/authors?%s", n.offset), nil)
		_, offset, _ := BindGetAuthorsRequest(s.req)
		s.Assert().Equal(n.expected, offset)
	}
}

var testsBooksId = []struct {
	input    string
	expected int
	err      error
}{
	{
		input:    "1",
		expected: 1,
		err:      nil,
	},
	{
		input:    "",
		expected: 0,
		err:      errors.New(erroId),
	},
	{
		input:    "-1",
		expected: 0,
		err:      errors.New(erroId),
	},
	{
		input:    "0",
		expected: 0,
		err:      errors.New(erroId),
	},
}

const erroId = "Erro teste id"

func (s *RequestBindingHandlerSuite) TestBooksIdPath() {
	for _, n := range testsBooksId {
		// arrange
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%s", n.input), nil)
		vars := map[string]string{
			"id": n.input,
		}

		s.req = mux.SetURLVars(s.req, vars)

		// act
		i, err := BindBookId(s.req, Path, erroId)

		// assert
		s.Assert().Equal(n.expected, i)
		s.Assert().Equal(n.err, err)
	}
}

func (s *RequestBindingHandlerSuite) TestBindGetBooksRequestNameQuery() {
	for _, n := range testsNome {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books?name=%s", n.nome), nil)
		request := BindGetBooksRequest(s.req)
		s.Assert().Equal(n.expected, request.Name)
	}
}

func (s *RequestBindingHandlerSuite) TestBindGetBooksRequestEditionQuery() {
	for _, n := range testsNome {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books?edition=%s", n.nome), nil)
		request := BindGetBooksRequest(s.req)
		s.Assert().Equal(n.expected, request.Edition)
	}
}

// se o publicationYear não for informado, deve ser usado 0
// se o publicationYear não for numero deve ser usado 0
// se o publicationYear for maior que o ano atual, o ano atual
// se o publicationYear for < 1500 deve ser usado 0
var testsPublicationYear = []struct {
	year     string
	expected int
}{
	{
		year:     "",
		expected: 0,
	},
	{
		year:     "notANumber&&",
		expected: 0,
	},
	{
		year:     fmt.Sprintf("%d", time.Now().Year()+1),
		expected: 0,
	},
	{
		year:     "1499",
		expected: 0,
	},
	{
		year:     fmt.Sprintf("%d", time.Now().Year()),
		expected: time.Now().Year(),
	},
}

func (s *RequestBindingHandlerSuite) TestBindGetBooksRequestPublicationYearQuery() {
	for _, n := range testsPublicationYear {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books?publication_year=%s", n.year), nil)
		request := BindGetBooksRequest(s.req)
		s.Assert().Equal(n.expected, request.PublicationYear)
	}
}

func (s *RequestBindingHandlerSuite) TestBindGetBooksRequestLimitQuery() {
	for _, n := range testsLimit {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books?%s", n.limit), nil)
		request := BindGetBooksRequest(s.req)
		s.Assert().Equal(n.expected, request.Limit)
	}
}

var testsValidPositiveInteger = []struct {
	value    string
	expected int
}{
	{
		value:    "",
		expected: 0,
	},
	{
		value:    "notANumber&&",
		expected: 0,
	},
	{
		value:    "10",
		expected: 10,
	},
	{
		value:    "-10",
		expected: 0,
	},
}

func (s *RequestBindingHandlerSuite) TestBindGetBooksRequestAuthorQuery() {
	for _, n := range testsValidPositiveInteger {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books?author=%s", n.value), nil)
		request := BindGetBooksRequest(s.req)
		s.Assert().Equal(n.expected, request.Author)
	}
}

func (s *RequestBindingHandlerSuite) TestBindCreateBookRequestBody() {

	body := domain.Book{
		Name:            "Some name",
		Edition:         "Some edition",
		PublicationYear: 2020,
	}
	json, _ := json.Marshal(body)
	s.req = httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(json))
	s.req.Header.Set("Content-Type", "application/json")

	request, err := BindCreateBookRequest(s.req)
	s.Assert().Equal(body, request)
	s.Assert().Nil(err)
}

func (s *RequestBindingHandlerSuite) TestBindCreateBookRequestBodyError() {

	s.req = httptest.NewRequest(http.MethodPost, "/books", nil)

	request, err := BindCreateBookRequest(s.req)
	s.Assert().Equal(domain.Book{}, request)
	s.Assert().Equal(errors.New("Invalid request payload"), err)
}

func TestBooksApiHandlerSuite(t *testing.T) {
	suite.Run(t, new(RequestBindingHandlerSuite))
}
