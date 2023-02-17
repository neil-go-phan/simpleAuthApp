package controllers

import (
	"errors"
	"example/web-service-gin/database"
	"example/web-service-gin/models"
	"regexp"

	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var validate *validator.Validate

type UserRepo struct {
	Db *gorm.DB
}

func New() *UserRepo {
	db := database.InitDb()
	// db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

func (repository *UserRepo) SignUp(c *gin.Context) {
	var user,checkUser models.User
	c.BindJSON(&user)
	err := validateSignUp(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong with input"})
		return
	}

	err = models.FindUser(repository.Db, &checkUser, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Bad request"})
		return
	}
	if checkUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Username is already taken"})
		return
	}
	err = models.CreateUser(repository.Db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Sign up success"})
}

func (repository *UserRepo) SignIn(c *gin.Context) {
	var user,checkUser models.User

	c.BindJSON(&user)
	err := validateSignIn(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong with input"})
		return
	}
	err = models.FindUser(repository.Db, &checkUser, user.Username)
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

func (repository *UserRepo) GetUsers(c *gin.Context) {
	var user []models.User
	err := models.GetUsers(repository.Db, &user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func validateSignUp(user *models.User) error {
	validate = validator.New()
	err := validate.Var(user.FullName, "required,min=8,max=50")
	if err != nil {
		return err
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
	match := checkRegexp(user.FullName, "full_name")
	if !match {
		return errors.New("full name must not contain special character")
	}
	match = checkRegexp(user.Username, "username")
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
	case "username":
		match, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", checkedString)
		return match
	case "full_name":
		match, _ := regexp.MatchString("^[a-zA-Z0-9_ ]*$", checkedString)
		return match
	}
	return false
}
