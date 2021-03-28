package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

type BooksApiHandlerSuite struct {
	suite.Suite

	ctx    context.Context
	router *mux.Router
	req    *http.Request
	res    *httptest.ResponseRecorder

	repo *database.BooksRepositoryMock
}

func (s *BooksApiHandlerSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = database.NewBooksRepositoryMock()
	s.res = httptest.NewRecorder()
	s.router = mux.NewRouter()
	NewBooksApi(s.router, s.repo)
}

func (s *BooksApiHandlerSuite) TestIfUrlInvalidShouldNotBe200() {
	s.req = httptest.NewRequest(http.MethodGet, "/bookss", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().NotEqual(http.StatusOK, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestIfUrlInvalidShouldReturn404() {
	s.req = httptest.NewRequest(http.MethodGet, "/bookss", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusNotFound, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestGetBooksShouldReturn200IfReturnedNil() {
	// arrange
	s.repo.
		On("GetAll", mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodGet, "/books", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	result := string(s.res.Body.String())

	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().JSONEq("[]", result)
}

func (s *BooksApiHandlerSuite) TestGetBooksShouldReturn200IfNotReturnedData() {

	// arrange
	s.repo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return([]domain.Book{})

	s.req = httptest.NewRequest(http.MethodGet, "/books", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	result := string(s.res.Body.String())

	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().JSONEq("[]", result)
}

func (s *BooksApiHandlerSuite) TestGetBooksShouldReturn200IfReturnedData() {

	// arrange
	books := []domain.Book{
		{
			Name:            "Book 1",
			Edition:         "1",
			PublicationYear: 2020,
		}, {
			Name:            "Book 2",
			Edition:         "2",
			PublicationYear: 2021,
		}}

	s.repo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return(books)

	s.req = httptest.NewRequest(http.MethodGet, "/books", nil)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	var result []domain.Book
	err := json.Unmarshal(s.res.Body.Bytes(), &result)
	if err != nil {
		log.Fatal(err)
	}

	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusOK, s.res.Code)
	s.Assert().Equal(books, result)
}

func (s *BooksApiHandlerSuite) TestGetBookShouldReturn404IfBookDoesNotExist() {
	// arrange
	s.repo.
		On("GetBook", mock.Anything).
		Return(nil, errors.New("Book not found"))

	s.req = httptest.NewRequest(http.MethodGet, "/books/1", nil)
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusNotFound, s.res.Code)
}

//se retornar alguma coisa deve ser 200
func (s *BooksApiHandlerSuite) TestGetBookShouldReturn200IfBookExists() {
	// arrange
	s.repo.
		On("GetBook", mock.Anything).
		Return(domain.Book{}, nil)

	s.req = httptest.NewRequest(http.MethodGet, "/books/1", nil)
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert

	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusOK, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestGetBookShouldReturn400IfIdIsInvalid() {

	// arrange
	s.repo.
		On("GetBook", mock.Anything).
		Return(domain.Book{}, nil)

	s.req = httptest.NewRequest(http.MethodGet, "/books/-1", nil)
	vars := map[string]string{
		"id": "-1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)

}

func (s *BooksApiHandlerSuite) TestPostBookShouldReturn200() {
	// arrange
	s.repo.
		On("Create", mock.Anything).
		Return(0, nil)

	var book = domain.Book{
		Name:            "Ronaldo",
		Edition:         "1",
		PublicationYear: 2020,
	}
	var jsonBytes, _ = json.Marshal(book)
	s.req = httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBytes))
	s.req.Header.Set("Content-Type", "application/json")

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusOK, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestPostBookShouldReturnIfBodyIsInvalid400() {
	// arrange
	s.repo.
		On("Create", mock.Anything).
		Return(0, nil)

	s.req = httptest.NewRequest(http.MethodPost, "/books", nil)
	s.req.Header.Set("Content-Type", "application/json")

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestPostBookShouldReturn400IfNotCreate() {
	// arrange
	s.repo.
		On("Create", mock.Anything).
		Return(0, errors.New("Some Error"))

	var book = domain.Book{
		Name:            "Ronaldo",
		Edition:         "1",
		PublicationYear: 2020,
	}
	var jsonBytes, _ = json.Marshal(book)
	s.req = httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBytes))
	s.req.Header.Set("Content-Type", "application/json")

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestPutBookShouldReturn200() {
	// arrange
	s.repo.
		On("Update", mock.Anything, mock.Anything).
		Return(nil)

	var book = domain.Book{
		Name:            "Ronaldo",
		Edition:         "1",
		PublicationYear: 2020,
	}
	var jsonBytes, _ = json.Marshal(book)
	s.req = httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewBuffer(jsonBytes))
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusOK, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestPutBookShouldReturn400IfHasNoBody() {
	// arrange
	s.repo.
		On("Update", mock.Anything, mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodPut, "/books/1", nil)
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestPutBookShouldReturn400IfIdIsLessThanZero() {
	// arrange
	s.repo.
		On("Update", mock.Anything, mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodPut, "/books/0", nil)
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "0",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestPutBookShouldReturn400UpdateFails() {
	// arrange
	s.repo.
		On("Update", mock.Anything, mock.Anything).
		Return(errors.New(""))

	var book = domain.Book{
		Name:            "Ronaldo",
		Edition:         "1",
		PublicationYear: 2020,
	}
	var jsonBytes, _ = json.Marshal(book)
	s.req = httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewBuffer(jsonBytes))
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestDeleteBookShouldReturn400IfIdIsLessThanZero() {
	// arrange
	s.repo.
		On("Delete", mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodDelete, "/books/0", nil)
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "0",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestDeleteBookShouldReturn400DeleteFails() {
	// arrange
	s.repo.
		On("Delete", mock.Anything).
		Return(errors.New("Some error"))

	s.req = httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *BooksApiHandlerSuite) TestDeleteBookShouldReturn200() {
	// arrange
	s.repo.
		On("Delete", mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	s.req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"id": "1",
	}

	s.req = mux.SetURLVars(s.req, vars)

	// act
	s.router.ServeHTTP(s.res, s.req)

	// assert
	s.Assert().Equal(http.StatusOK, s.res.Code)
}

func TestBooksApiHandlerSuite(t *testing.T) {
	suite.Run(t, new(BooksApiHandlerSuite))
}
