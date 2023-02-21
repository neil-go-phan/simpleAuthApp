package controllers

import (
	"example/web-service-gin/models"
	"testing"
)

type checkRegexpTestCase struct {
	testName  string
	str       string
	checkType string
	isPass    bool
}

var checkRegexpTestCases = []checkRegexpTestCase{
	{testName: "pass", str: "Golang 2023", checkType: "usernameAndPassword", isPass: false},
	{testName: "username contain &*", str: "Golang&*", checkType: "usernameAndPassword", isPass: false},
	{testName: "username contain __", str: "Golang___2023", checkType: "usernameAndPassword", isPass: true},
	{testName: "full name regex pass", str: "Golang 2023", checkType: "full_name", isPass: true},
	{testName: "full name contain %^", str: "Golang%^2023", checkType: "full_name", isPass: false},
	{testName: "full name contain __", str: "Golang__2023", checkType: "full_name", isPass: true},
}

func assertCheckRegexp(t *testing.T, str string, checkType string, isPass bool) {
	want := isPass
	got := checkRegexp(str, checkType)
	if got != want {
		t.Errorf("%s with checkType = '%s' is supose to %v", str, checkType, want)
	}
}
func TestCheckRegexp(t *testing.T) {
	for _, c := range checkRegexpTestCases {
		t.Run(c.testName, func(t *testing.T) {
			assertCheckRegexp(t, c.str, c.checkType, c.isPass)
		})
	}
}

type validateUsernameAndPasswordAndSignInTestCase struct {
	user     models.User
	isPass   bool
	testName string
}

var validateUsernameAndPasswordAndSignInTestCases = []validateUsernameAndPasswordAndSignInTestCase{
	{models.User{Username: "golang123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, true, "pass all"},
	{models.User{Username: "gol", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username too short"},
	{models.User{Username: "golang123", Password: "123456789"}, false, "password is not hash"},
	{models.User{Username: "golang 123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username contain space"},
	{models.User{Username: "golang_123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, true, "username contain __"},
	{models.User{Username: "golang123&^%", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username contain &^%"},
	{models.User{Username: "golang1231234120941200912801123213124123123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username too long"},
}

func assertValidateUsernameAndPassword(t *testing.T, user models.User, isPass bool) {
	want := isPass
	var result bool
	got := validateUsernameAndPassword(&user)
	if got != nil {
		result = false
	} else {
		result = true
	}
	if result != want {
		t.Errorf("username = %s\npassword = %s\nis suppose to be %v", user.Username, user.Password, want)
	}
}
func TestValidateUsernameAndPassword(t *testing.T) {
	for _, c := range validateUsernameAndPasswordAndSignInTestCases {
		t.Run(c.testName, func(t *testing.T) {
			assertValidateUsernameAndPassword(t, c.user, c.isPass)
		})
	}
}

func assertValidateSignIn(t *testing.T, user models.User, isPass bool) {
	want := isPass
	var result bool
	got := validateSignIn(&user)
	if got != nil {
		result = false
	} else {
		result = true
	}
	if result != want {
		t.Errorf("username = %s\npassword = %s\nis suppose to be %v", user.Username, user.Password, want)
	}
}
func TestValidateSignIn(t *testing.T) {
	for _, c := range validateUsernameAndPasswordAndSignInTestCases {
		t.Run(c.testName, func(t *testing.T) {
			assertValidateSignIn(t, c.user, c.isPass)
		})
	}
}

type validateSignUpTestCase struct {
	user     models.User
	isPass   bool
	testName string
}

var validateSignUpTestCases = []validateSignUpTestCase{
	{models.User{FullName: "Golden Owl 2023", Username: "golang123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, true, "pass"},
	{models.User{FullName: "Golden Owl 2023", Username: "gol", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username too short"},
	{models.User{FullName: "Golden Owl 2023", Username: "golang123", Password: "123456789"}, false, "password is not hash"},
	{models.User{FullName: "Golden Owl 2023", Username: "golang 123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username contain space"},
	{models.User{FullName: "Golden Owl 2023", Username: "golang_123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, true, "username contain __"},
	{models.User{FullName: "Golden Owl 2023", Username: "golang123&^%", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username contain &^%"},
	{models.User{FullName: "Golden Owl 2023", Username: "golang1231234120941200912801123213124123123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "username too long"},
	{models.User{FullName: "Golden__Owl_2023", Username: "golang123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, true, "full name contain __"},
	{models.User{FullName: "Gold", Username: "golang123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "full name too short"},
	{models.User{FullName: "Golden Owl&^%", Username: "golang123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "full name contain &^%"},
	{models.User{FullName: "Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl Golden Owl", Username: "golang123", Password: "81dc9bdb52d04dc20036dbd8313ed055"}, false, "full name too long"},
}

func assertValidateSignUp(t *testing.T, user models.User, isPass bool) {
	want := isPass
	var result bool
	got := validateSignUp(&user)
	if got != nil {
		result = false
	} else {
		result = true
	}
	if result != want {
		t.Errorf("username = %s\npassword = %s\nis suppose to be %v", user.Username, user.Password, want)
	}
}
func TestValidateSignUp(t *testing.T) {
	for _, c := range validateSignUpTestCases {
		t.Run(c.testName, func(t *testing.T) {
			assertValidateSignUp(t, c.user, c.isPass)
		})

	}
}

