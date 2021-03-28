package database

import (
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
	Create(book domain.Book) (uint, error)
	Update(id int, book domain.Book) error
	Delete(id int) error
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
	users := []domain.Book{}

	db := i.manager.GetDB()

	if len(r.Name) > 0 {
		db.Where("name = ?", r.Name)
	}

	if len(r.Edition) > 0 {
		db.Where("edition = ?", r.Edition)
	}

	if r.PublicationYear > 0 {
		db.Where("publication_year = ?", r.PublicationYear)
	}

	db.Take(r.Limit)
	db.Offset(r.Offset)
	db.Find(&users)
	return users
}

func (i *booksRepository) GetBook(id int) (domain.Book, error) {
	b := domain.Book{}
	err := i.manager.GetDB().First(&b).Error
	return b, err
}

func (i *booksRepository) Create(b domain.Book) (uint, error) {
	book := domain.Book{
		Name:            b.Name,
		Edition:         b.Edition,
		PublicationYear: b.PublicationYear,
	}

	err := i.manager.GetDB().Create(&book).Error
	return book.ID, err
}

func (i *booksRepository) Update(id int, b domain.Book) error {
	book := domain.Book{}
	i.manager.GetDB().First(&book, id)

	book.Name = b.Name
	book.Edition = b.Edition
	book.PublicationYear = b.PublicationYear

	return i.manager.GetDB().Save(&book).Error
}

func (i *booksRepository) Delete(id int) error {
	manager := i.manager.GetDB()
	return manager.Delete(&domain.Book{}, id).Error
}
