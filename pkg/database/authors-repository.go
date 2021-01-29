package database

import (
	"fmt"

	"github.com/jedielson/bookstore/pkg/domain"
)

type AuthorsRepository interface {
	GetAll(name string, limit int, offset int) []domain.Author
}

type authorsRepository struct {
	manager DBManager
}

func NewAuthorsRepository(m DBManager) AuthorsRepository {
	return &authorsRepository{
		manager: m,
	}
}

func (a *authorsRepository) GetAll(name string, limit int, offset int) []domain.Author {
	var records []domain.Author
	var db = a.manager.GetDB()

	if len(name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	err := db.Limit(limit).Offset(offset).Find(&records)

	if err.Error != nil {
		return nil
	}

	return records
}
