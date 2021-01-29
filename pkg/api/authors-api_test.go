package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthorsApiHandlerSuite struct {
	suite.Suite

	ctx context.Context

	req *http.Request
	res *httptest.ResponseRecorder

	repo *database.AuthorsRepositoryMock
}

func (s *AuthorsApiHandlerSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = database.NewAuthorsRepositoryMock()
	s.res = httptest.NewRecorder()
}

func (s *AuthorsApiHandlerSuite) TestShouldReturn200IfReturnedNil() {

	// arrange
	s.repo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodGet, "/authors", nil)

	// act
	GetAuthors(s.repo)(s.res, s.req)

	// assert
	result := string(s.res.Body.String())

	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().JSONEq("[]", result)
}

func (s *AuthorsApiHandlerSuite) TestShouldReturn200IfNotReturnedData() {

	// arrange
	s.repo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return([]domain.Author{})

	s.req = httptest.NewRequest(http.MethodGet, "/authors", nil)

	// act
	GetAuthors(s.repo)(s.res, s.req)

	// assert
	result := string(s.res.Body.String())

	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().JSONEq("[]", result)
}

func (s *AuthorsApiHandlerSuite) TestShouldReturn200IfReturnedData() {

	// arrange
	authors := []domain.Author{
		{
			Name: "Teste",
		}, {
			Name: "Teste2",
		}}

	s.repo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return(authors)

	s.req = httptest.NewRequest(http.MethodGet, "/authors", nil)

	// act
	GetAuthors(s.repo)(s.res, s.req)

	// assert
	var result []domain.Author
	err := json.Unmarshal(s.res.Body.Bytes(), &result)
	if err != nil {
		log.Fatal(err)
	}

	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().Equal(authors, result)
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
		nome:     "name=myname",
		expected: "myname",
	},
}

func (s *AuthorsApiHandlerSuite) TestNameQuery() {
	for _, n := range testsNome {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/authors?%s", n.nome), nil)
		nome, _, _ := parseUrl(s.req)
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

func (s *AuthorsApiHandlerSuite) TestLimitQuery() {
	for _, n := range testsLimit {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/authors?%s", n.limit), nil)
		_, _, limit := parseUrl(s.req)
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

func (s *AuthorsApiHandlerSuite) TestOffsetQuery() {
	for _, n := range testsOffset {
		s.req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/authors?%s", n.offset), nil)
		_, offset, _ := parseUrl(s.req)
		s.Assert().Equal(n.expected, offset)
	}
}

func TestAuthorsApiHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthorsApiHandlerSuite))
}
