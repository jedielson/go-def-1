package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jedielson/bookstore/pkg/database"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BooksApiHandlerSuite struct {
	suite.Suite

	ctx context.Context
	req *http.Request
	res *httptest.ResponseRecorder

	repo *database.BooksRepositoryMock
}

func (s *BooksApiHandlerSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = database.NewBooksRepositoryMock()
	s.res = httptest.NewRecorder()
}

func (s *BooksApiHandlerSuite) TestShouldReturn200IfReturnedNil() {
	// arrange
	s.repo.
		On("GetAll", mock.Anything).
		Return(nil)

	s.req = httptest.NewRequest(http.MethodGet, "/books", nil)

	// act
	GetBooks(s.repo)(s.res, s.req)

	// assert
	//result := string(s.res.Body.String())

	s.repo.AssertExpectations(s.T())
	s.Assert().Equal(http.StatusOK, s.res.Code)
	//s.Assert().JSONEq("[]", result)
}

// func (s *BooksApiHandlerSuite) TestShouldReturn200IfNotReturnedData() {

// 	// arrange
// 	s.repo.
// 		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
// 		Return([]domain.Book{})

// 	s.req = httptest.NewRequest(http.MethodGet, "/books", nil)

// 	// act
// 	GetBooks(s.repo)(s.res, s.req)

// 	// assert
// 	result := string(s.res.Body.String())

// 	s.Assert().Equal(http.StatusOK, s.res.Code)
// 	s.Assert().JSONEq("[]", result)
// }

// func (s *BooksApiHandlerSuite) TestShouldReturn200IfReturnedData() {

// 	// arrange
// 	books := []domain.Book{
// 		{
// 			Name:            "Book 1",
// 			Edition:         "1",
// 			PublicationYear: 2020,
// 		}, {
// 			Name:            "Book 2",
// 			Edition:         "2",
// 			PublicationYear: 2021,
// 		}}

// 	s.repo.
// 		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
// 		Return(books)

// 	s.req = httptest.NewRequest(http.MethodGet, "/books", nil)

// 	// act
// 	GetBooks(s.repo)(s.res, s.req)

// 	// assert
// 	var result []domain.Book
// 	err := json.Unmarshal(s.res.Body.Bytes(), &result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	s.Assert().Equal(http.StatusOK, s.res.Code)
// 	s.Assert().Equal(books, result)
// }

func TestBooksApiHandlerSuite(t *testing.T) {
	suite.Run(t, new(BooksApiHandlerSuite))
}
