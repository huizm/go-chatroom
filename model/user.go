package model

import "gorm.io/gorm"

type User struct {
	Username string `gorm:"uniqueIndex,notNull" json:"username,omitempty"`
	Password string `gorm:"notNull" json:"password,omitempty"`

	gorm.Model
}

func (*User) TableName() string {
	return "users"
}

func SearchUser(data *User) ([]*User, error) {
	tx := DBEngine

	var users []*User
	err := tx.Where(data).Find(&users).Error
	return users, err
}

func CreateUser(data *User) error {
	return DBEngine.Create(data).Error
}
