package repository

import (
	"golang-crud-ddd/domain"
	models "golang-crud-ddd/model"

	"github.com/beego/beego/v2/client/orm"
)

type UserRepo struct {
	// Add any necessary fields here, such as a database connection
	o orm.Ormer
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		o: orm.NewOrm(),
	}
}

// GetUserByID retrieves a user by ID
func (repo *UserRepo) GetUserByID(id int) (*domain.User, error) {
	user, err := models.GetUserByID(id, repo.o)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByEmail ...
func (repo *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	user, err := models.GetUserByEmail(email, repo.o)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// CreateUser creates a new user
func (repo *UserRepo) CreateUser(user *domain.User) (int64, error) {
	modelUser := &models.Users{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}
	id, err := models.CreateUser(modelUser, repo.o)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateUserByID updates a user by ID
func (repo *UserRepo) UpdateUserByID(id int, user *domain.User) error {
	modelUser := &models.Users{
		ID:       id,
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}
	err := models.UpdateUserByID(id, modelUser, repo.o)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserByID deletes a user by ID
func (repo *UserRepo) DeleteUserByID(id int) error {
	err := models.DeleteUserByID(id, repo.o)
	if err != nil {
		return err
	}
	return nil
}

// ChangePassword ...
func (repo *UserRepo) ChangePassword(id int, newPassword string) error {
	err := models.ChangePassword(id, newPassword, repo.o)
	if err != nil {
		return err
	}
	return nil
}