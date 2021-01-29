package database

import (
	"fmt"

	"github.com/jedielson/bookstore/pkg/domain"
)

type GetAllRequest struct {
	name             string
	publication_year int
	edition          string
	author           int
	limit            int
	offset           int
}

type BooksRepository interface {
	GetAll(r GetAllRequest) []domain.Book
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
			Name:            r.name,
			PublicationYear: r.publication_year,
			Edition:         r.edition,
			Authors: []*domain.Author{
				{
					Name: fmt.Sprintf("%d", r.author),
				},
			},
		},
	}
}
