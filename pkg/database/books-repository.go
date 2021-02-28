package database

import (
	"errors"
	"fmt"

	"github.com/jedielson/bookstore/pkg/domain"
)

type GetAllRequest struct {
	Name            string
	PublicationYear int
	Edition         string
	Author          int
	Limit           int
	Offset          int
}

type BooksRepository interface {
	GetAll(r GetAllRequest) []domain.Book
	GetBook(id int) (domain.Book, error)
}

type booksRepository struct {
	manager DBManager
}

func NewBooksRepository(d DBManager) BooksRepository {
	return &booksRepository{
		manager: d,
	}
}

func (i *booksRepository) GetAll(r GetAllRequest) []domain.Book {
	return []domain.Book{
		{
			Name:            r.Name,
			PublicationYear: r.PublicationYear,
			Edition:         r.Edition,
			Authors: []*domain.Author{
				{
					Name: fmt.Sprintf("%d", r.Author),
				},
			},
		},
	}
}

func (i *booksRepository) GetBook(id int) (domain.Book, error) {
	return domain.Book{}, errors.New("Method not implemented")
}
