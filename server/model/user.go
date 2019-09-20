package model

import (
	"github.com/tidyoux/router/common/db"
	"github.com/tidyoux/router/common/types"
)

type User struct {
	Model

	Name     string `gorm:"size:32;unique_index"`
	Password string `gorm:"size:64"`
	Detail   string `gorm:"type:text"`
	Status   int8   `gorm:"type:tinyint;index"`
}

func NewUser(name, password string) *User {
	return &User{
		Name:     name,
		Password: password,
	}
}

func (*User) TableName() string { return "user" }

func (u *User) Insert() error {
	return db.Default().Create(u).Error
}

func (u *User) UpdatePassword(password string) error {
	return u.update(M{
		"password": password,
	})
}

func (u *User) UpdateDetail(detail string) error {
	return u.update(M{
		"detail": detail,
	})
}

func (u *User) Enable() error {
	return u.update(M{
		"status": types.UserEnabled,
	})
}

func (u *User) Disable() error {
	return u.update(M{
		"status": types.UserDisabled,
	})
}

func (u *User) IsAdmin() bool {
	return u.Name == types.AdminUsername
}

func (u *User) update(values M) error {
	return db.Default().Model(u).Updates(values).Error
}

func FindAllUsers() ([]*User, error) {
	var users []*User
	err := db.Default().Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func FindUserByID(id uint64) (*User, error) {
	var user User
	err := db.Default().First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByName(name string) (*User, error) {
	var user User
	err := db.Default().First(&user, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
