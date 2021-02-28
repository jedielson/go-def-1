package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthorsApiHandlerSuite struct {
	suite.Suite

	ctx    context.Context
	router *mux.Router

	req *http.Request
	res *httptest.ResponseRecorder

	repo *database.AuthorsRepositoryMock
}

func (s *AuthorsApiHandlerSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = database.NewAuthorsRepositoryMock()
	s.res = httptest.NewRecorder()
	s.router = mux.NewRouter()
	NewAuthorsApi(s.router, s.repo)
}

func (s *AuthorsApiHandlerSuite) TestIfUrlInvalidShouldReturn404() {
	s.req = httptest.NewRequest(http.MethodGet, "/autors", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusNotFound, s.res.Code)
}

func (s *AuthorsApiHandlerSuite) TestShouldReturn200IfReturnedNil() {

	// arrange
	s.repo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodGet, "/authors", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	result := string(s.res.Body.String())

	s.repo.AssertExpectations(s.T())
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
	s.router.ServeHTTP(s.res, s.req)

	// assert
	result := string(s.res.Body.String())

	s.repo.AssertExpectations(s.T())
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
	s.router.ServeHTTP(s.res, s.req)

	// assert
	var result []domain.Author
	err := json.Unmarshal(s.res.Body.Bytes(), &result)
	if err != nil {
		log.Fatal(err)
	}

	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().Equal(authors, result)
}

func TestAuthorsApiHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthorsApiHandlerSuite))
}
