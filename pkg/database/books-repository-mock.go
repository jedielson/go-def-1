package database

import (
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type BooksRepositoryMock struct {
	mock.Mock
}

func NewBooksRepositoryMock() *BooksRepositoryMock {
	return &BooksRepositoryMock{}
}

func (m *BooksRepositoryMock) GetAll(r GetAllRequest) []domain.Book {
	args := m.Called(r)
	bb, ok := args.Get(0).([]domain.Book)

	if !ok {
		return nil
	}

	return bb
}
