package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID        int       `orm:"column(id);auto"`
	Name      string    `orm:"column(name);size(100)"`
	Email     string    `orm:"column(email);size(150);unique"` // unique index
	Phone     string    `orm:"column(phone);size(20);null"`    // add unique if required: ;unique
	Password  string    `orm:"column(password);"`
	Role      string    `orm:"column(role);size(20)"` // rider|driver|admin (enforce in app or via CHECK in SQL)
	Avatar    string    `orm:"column(avatar);size(255);null"`
	IsActive  bool      `orm:"column(is_active);default(true)"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)"` // set once
	UpdatedAt time.Time `orm:"column(updated_at);auto_now;type(datetime)"`     // set on every save
}

func (u *Users) TableName() string { return "users" }

func init() {
	orm.RegisterModel(new(Users))
}

// Queries
func GetUserByID(id int, o orm.Ormer) (*Users, error) {
	var u Users
	if err := o.QueryTable(new(Users)).Filter("id", id).One(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(email string, o orm.Ormer) (*Users, error) {
	var u Users
	if err := o.QueryTable(new(Users)).Filter("email", email).One(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// Create
func CreateUser(u *Users, o orm.Ormer) (int64, error) {
	if u.Password != "" {
		enc, err := encryptPassword(u.Password)
		if err != nil {
			return 0, err
		}
		u.Password = enc
	}
	return o.Insert(u)
}

// Update: pass explicit fields to avoid overwriting unintended columns
func UpdateUserByID(id int, u *Users, o orm.Ormer, fields ...string) error {
	u.ID = id
	_, err := o.Update(u, fields...)
	return err
}

// Delete
func DeleteUserByID(id int, o orm.Ormer) error {
	_, err := o.Delete(&Users{ID: id})
	return err
}

// Password helpers
func encryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// Change password updates UpdatedAt automatically due to auto_now
func ChangePassword(id int, newPassword string, o orm.Ormer) error {
	enc, err := encryptPassword(newPassword)
	if err != nil {
		return err
	}
	u := &Users{ID: id, Password: enc}
	_, err = o.Update(u, "Password") // UpdatedAt auto_now will refresh on save
	return err
}
