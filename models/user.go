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

type Repository interface {
	CreateUser(userInput *User) (*User, error)
	GetUsers() (user *[]User,err error)
	FindUser(username string) (u *User, err error) 
}
type Repo struct {
	DB *gorm.DB
}
func (p *Repo)CreateUser(userInput *User) (*User, error) {
	user := User{
		FullName: userInput.FullName,
		Username: userInput.Username,
		Password: userInput.Password,
		AuthorityType: "customer",
	}
	err := p.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return &user ,nil
}

func (p *Repo)GetUsers() (user *[]User,err error) {
	users := make([]User, 10)
	err = p.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (p *Repo)FindUser(username string) (u *User, err error) {
	user := new(User)
	err = p.DB.Select("username", "password").Where(map[string]interface{}{"username": username}).Find(&user).Error
	if err!= nil {
		return nil, err
	}
	return user, nil
}

func CreateRepository(db *gorm.DB) Repository {
	return &Repo{
		DB: db,
	}
}