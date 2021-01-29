package database

import (
	"github.com/jedielson/bookstore/pkg/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(domain.Author{}, domain.Book{})

	if err != nil {
		panic(err)
	}
}
