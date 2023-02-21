package controllers

import (
	"errors"
	"example/web-service-gin/models"
	"regexp"

	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

var validate *validator.Validate

// type Repository interface {
// 	SignUp(c *gin.Context) ()
// 	SignIn(c *gin.Context) ()
// 	GetUsers(c *gin.Context) ()
// }

type Server struct {
	repo models.Repository
}

func NewServer(repo models.Repository) *Server{
	server := &Server{
		repo: repo,
	}
	return server
}

func (server *Server)SignUp(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	err := validateSignUp(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong with input"})
		return
	}

	checkUser, err := server.repo.FindUser(user.Username)
	// models.Repository.FindUser( &checkUser, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Bad request"})
		return
	}
	if checkUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Username is already taken"})
		return
	}
	_, err = server.repo.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Sign up success"})
}

func (server *Server)SignIn(c *gin.Context) {
	var user models.User

	c.BindJSON(&user)
	err := validateSignIn(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong with input"})
		return
	}
	checkUser, err := server.repo.FindUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Bad request"})
		return
	}
	if user.Username != checkUser.Username || user.Password != checkUser.Password {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Username or password is incorrect"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Sign in success"})
}

func (server *Server)GetUsers(c *gin.Context) {
	// var user []models.User
	users, err := server.repo.GetUsers()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func validateSignUp(user *models.User) error {
	validate = validator.New()
	err := validate.Var(user.FullName, "required,min=8,max=50")
	if err != nil {
		return err
	}
	checkRegexFullName := checkRegexp(user.FullName, "full_name")
	if !checkRegexFullName {
		return errors.New("full name must not contain special character")
	}
	err = validateUsernameAndPassword(user)
	if err != nil {
		return err
	}
	return nil
}

func validateSignIn(user *models.User) error {
	err := validateUsernameAndPassword(user)
	if err != nil {
		return err
	}
	return nil
}

func validateUsernameAndPassword(user *models.User) error {
	validate = validator.New()
	match := checkRegexp(user.Password, "usernameAndPassword")
	if !match {
		return errors.New("password must not contain special character")
	}
	match = checkRegexp(user.Username, "usernameAndPassword")
	if !match {
		return errors.New("username must not contain special character")
	}
	err := validate.Struct(user)
	if err != nil {
		return err
	}
	return nil
}

func checkRegexp(checkedString string, checkType string) bool {
	switch checkType {
	case "usernameAndPassword":
		match, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", checkedString)
		return match
	case "full_name":
		match, _ := regexp.MatchString("^[a-zA-Z0-9_ ]*$", checkedString)
		return match
	}
	return false
}
