package domain

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name      string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Books     []*Book        `gorm:"many2many:author_books;"`
}

type Book struct {
	gorm.Model
	Name            string `gorm:"size:255"`
	Edition         string `gorm:"size:255"`
	PublicationYear int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Authors         []*Author      `gorm:"many2many:author_books;"`
}
