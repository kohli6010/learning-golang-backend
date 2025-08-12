package models

import (
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type Users struct {
	ID        int    `orm:"column(id);auto"`
	Name      string `orm:"column(name);"`
	Password  string `orm:"column(password);"`
	Email     string `orm:"column(email);"`
	CreatedAt string `orm:"column(created_at);"`
	UpdatedAt string `orm:"column(updated_at);"`
}

func init() {
	orm.RegisterModel(new(Users))
}

// GetUserByID ... retrieves a user by ID
func GetUserByID(id int, o orm.Ormer) (*Users, error) {
	var user Users
	err := o.QueryTable(new(Users)).Filter("id", id).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail ...
func GetUserByEmail(email string, o orm.Ormer) (*Users, error) {
	var user Users
	err := o.QueryTable(new(Users)).Filter("email", email).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser ... creates a new user
func CreateUser(user *Users, o orm.Ormer) (int64, error) {
	// Insert the user into the database
	if user.Password != "" {
		encryptedPassword, err := encrpytPassword(user.Password)
		if err != nil {
			return 0, err
		}
		user.Password = encryptedPassword
	}

	id, err := o.Insert(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateUserByID ... updates a user by ID
func UpdateUserByID(id int, user *Users, o orm.Ormer) error {
	// Update the user in the database
	user.ID = id
	_, err := o.Update(user)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserByID ... deletes a user by ID
func DeleteUserByID(id int, o orm.Ormer) error {
	// Delete the user from the database
	_, err := o.Delete(&Users{ID: id})
	if err != nil {
		return err
	}
	return nil
}

// encrpytPassword ...
func encrpytPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ChangePassword ...
func ChangePassword(id int, newPassword string, o orm.Ormer) error {
	// Encrypt the new password
	encryptedPassword, err := encrpytPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the user's password in the database
	user := &Users{ID: id, Password: encryptedPassword}
	_, err = o.Update(user, "Password")
	if err != nil {
		return err
	}

	return nil
}