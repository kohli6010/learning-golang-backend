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
		Role:      user.Role,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		Password:  user.Password, // Ensure password is handled securely
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
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
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateUser creates a new user
func (repo *UserRepo) CreateUser(user *domain.User) (int64, error) {
	modelUser := &models.Users{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
		Role:     user.Role,
	}
	id, err := models.CreateUser(modelUser, repo.o)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateUserByID updates a user by ID
func (repo *UserRepo) UpdateUserByID(id int, user *domain.User) error {
	modelUser := &models.Users{}
	if user.Name != "" {
		modelUser.Name = user.Name
	}
	if user.Email != "" {
		modelUser.Email = user.Email
	}
	if user.Phone != "" {
		modelUser.Phone = user.Phone
	}
	if user.Role != "" {
		modelUser.Role = user.Role
	}
	if user.Avatar != "" {
		modelUser.Avatar = user.Avatar
	}
	if user.IsActive {
		modelUser.IsActive = user.IsActive
	}
	modelUser.Password = user.Password // Assuming password is already encrypted
	modelUser.IsActive = true          // Default to true if not specified
	
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
