//go:generate mockery --name=IUserRepository --output=../mocks --outpkg=mocks --case=underscore

package repository

import "golang-crud-ddd/domain"

type IUserRepository interface {
	GetUserByID(id int) (*domain.User, error)
	CreateUser(user *domain.User) (int64, error)
	UpdateUserByID(id int, user *domain.User) error
	DeleteUserByID(id int) error
	GetUserByEmail(email string) (*domain.User, error)
	ChangePassword(id int, newPassword string) error
}