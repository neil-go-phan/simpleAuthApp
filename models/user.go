package models

import (
	"gorm.io/gorm"
)

type User struct {
	FullName string `json:"full_name"`
	Username string	`json:"username" validate:"required,min=8,max=16"`
	Password string `json:"password" validate:"required,md5"`
	AuthorityType string `json:"authority_type"`
}

// struct tag 

//create a user
func CreateUser(db *gorm.DB, User *User) (err error) {
	User.AuthorityType = "customer"
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get users
func GetUsers(db *gorm.DB, User *[]User) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

func FindUser(db *gorm.DB, User *User, username string) (err error) {
	err = db.Select("username", "password").Where(map[string]interface{}{"username": username}).Find(&User).Error

	if err != nil {
		return err
	}
	return nil
}
