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

func (m *BooksRepositoryMock) GetBook(id int) (domain.Book, error) {
	args := m.Called(id)
	bb, ok := args.Get(0).(domain.Book)

	if !ok {
		bb = domain.Book{}
	}

	return bb, args.Error(1)
}

func (m *BooksRepositoryMock) Create(book domain.Book) (uint, error) {
	args := m.Called(book)
	bb, ok := args.Get(0).(uint)

	if !ok {
		bb = 0
	}

	return bb, args.Error(1)
}

func (m *BooksRepositoryMock) Update(id int, book domain.Book) error {
	args := m.Called(id, book)
	return args.Error(0)
}

func (m *BooksRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
