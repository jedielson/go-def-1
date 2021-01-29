package database

import (
	"github.com/jedielson/bookstore/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type AuthorsRepositoryMock struct {
	mock.Mock
}

func NewAuthorsRepositoryMock() *AuthorsRepositoryMock {
	return &AuthorsRepositoryMock{}
}

func (m *AuthorsRepositoryMock) GetAll(name string, limit int, offset int) []domain.Author {
	args := m.Called(name, limit, offset)
	bb, ok := args.Get(0).([]domain.Author)

	if !ok {
		return nil
	}

	return bb
}
