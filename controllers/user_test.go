package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"example/web-service-gin/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type checkRegexpTestCase struct {
	testName  string
	str       string
	checkType string
	expected  bool
}

type SuiteController struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository models.Repository
	server     Server
}

type validateUsernameAndPasswordAndSignInTestCase struct {
	user     models.User
	expected bool
	testName string
}

type validateSignUpTestCase struct {
	user     models.User
	expected bool
	testName string
}

type requestTestCase struct {
	testName       string
	expected       bool
	input          models.User
	databaseReturn []models.User
	responseJSONs  map[string]interface{}
	statusCode     int
}

var checkRegexpTestCases = []checkRegexpTestCase{
	{testName: "pass", str: "Golang 2023", checkType: "usernameAndPassword", expected: false},
	{testName: "username contain &*", str: "Golang&*", checkType: "usernameAndPassword", expected: false},
	{testName: "username contain __", str: "Golang___2023", checkType: "usernameAndPassword", expected: true},
	{testName: "full name regex pass", str: "Golang 2023", checkType: "full_name", expected: true},
	{testName: "full name contain %^", str: "Golang%^2023", checkType: "full_name", expected: false},
	{testName: "full name contain __", str: "Golang__2023", checkType: "full_name", expected: true},
}

func assertCheckRegexp(t *testing.T, str string, checkType string, expected bool) {
	want := expected
	got := checkRegexp(str, checkType)
	if got != want {
		t.Errorf("%s with checkType = '%s' is supose to %v", str, checkType, want)
	}
}
func TestCheckRegexp(t *testing.T) {
	for _, c := range checkRegexpTestCases {
		t.Run(c.testName, func(t *testing.T) {
			assertCheckRegexp(t, c.str, c.checkType, c.expected)
		})
	}
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

func assertValidateUsernameAndPassword(t *testing.T, user models.User, expected bool) {
	want := expected
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
			assertValidateUsernameAndPassword(t, c.user, c.expected)
		})
	}
}

func assertValidateSignIn(t *testing.T, user models.User, expected bool) {
	want := expected
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
			assertValidateSignIn(t, c.user, c.expected)
		})
	}
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

func assertValidateSignUp(t *testing.T, user models.User, expected bool) {
	want := expected
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
			assertValidateSignUp(t, c.user, c.expected)
		})
	}
}

// mock database
func (s *SuiteController) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err)

	conn, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	s.DB = conn
	require.NoError(s.T(), err)

	s.repository = models.CreateRepository(s.DB)
	s.server = *NewServer(s.repository)
}
func (s *SuiteController) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(SuiteController))
}

var getUsersTestCases = []requestTestCase{
	{
		testName: "success to get users",
		expected: true,
		databaseReturn: []models.User{
			{FullName: "admin", Username: "admin1234", Password: "81dc9bdb52d04dc20036dbd8313ed055", AuthorityType: "admin"},
			{FullName: "customer", Username: "customer1234", Password: "81dc9bdb52d04dc20036dbd8313ed055", AuthorityType: "customer"},
			{FullName: "employee", Username: "employee1234", Password: "81dc9bdb52d04dc20036dbd8313ed055", AuthorityType: "employee"},
		},
		responseJSONs: nil,
		statusCode:    http.StatusOK,
	},
	{
		testName:       "fail to get users",
		expected:       false,
		databaseReturn: nil,
		responseJSONs: map[string]interface{}{
			"success": false,
			"message": "Bad request",
		},
		statusCode: http.StatusBadRequest,
	},
}

func (s *SuiteController) TestGetUser() {
	for _, c := range getUsersTestCases {
		s.T().Run(c.testName, func(t *testing.T) {
			assertTestGetUser(s, c)
		})
	}
}

func assertTestGetUser(s *SuiteController, testCase requestTestCase) {
	// config result
	var want any
	if testCase.expected == true {
		want = testCase.databaseReturn

	} else {
		want = testCase.responseJSONs
	}
	// mock database
	if testCase.expected == true {
		rows := sqlmock.NewRows([]string{"full_name", "username", "password", "authority_type"})
		for _, r := range testCase.databaseReturn {
			rows.AddRow(r.FullName, r.Username, r.Password, r.AuthorityType)
		}

		s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnRows(rows)

	} else {
		s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnError(errors.New("fail to get users"))
	}

	// mock a request
	r := gin.Default()
	r.GET("auth/users", s.server.GetUsers)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/auth/users", nil)
	assert.Nil(s.T(), err)
	r.ServeHTTP(w, req)
	assert.Equal(s.T(), testCase.statusCode, w.Code)
	if testCase.expected == true {
		// parse reponse
		var got []models.User
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), want, got)
	} else {
		// parse reponse
		var got map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), want, got)
	}
}

var signInTestCases = []requestTestCase{
	{
		testName: "success to sign in",
		expected: true,
		input:    models.User{Username: "admin1234", Password: "81dc9bdb52d04dc20036dbd8313ed055"},
		databaseReturn: []models.User{
			{Username: "admin1234", Password: "81dc9bdb52d04dc20036dbd8313ed055"},
		},
		responseJSONs: map[string]interface{}{
			"success": true,
			"message": "Sign in success",
		},
		statusCode: http.StatusOK,
	},
	{
		testName:       "username incorrect",
		expected:       false,
		databaseReturn: nil,
		input:          models.User{Username: "admin2321", Password: "81dc9bdb52d04dc20036dbd8313ed055"},
		responseJSONs: map[string]interface{}{
			"success": false,
			"message": "Username or password is incorrect",
		},
		statusCode: http.StatusBadRequest,
	},
	{
		testName:       "password incorrect",
		expected:       false,
		databaseReturn: nil,
		input:          models.User{Username: "admin1234", Password: "844c81785ce3fe7889cfae15bb570820"},
		responseJSONs: map[string]interface{}{
			"success": false,
			"message": "Username or password is incorrect",
		},
		statusCode: http.StatusBadRequest,
	},
}

func (s *SuiteController) TestSignIn() {
	for _, c := range signInTestCases {
		s.T().Run(c.testName, func(t *testing.T) {
			assertTestSignIn(s, c)
		})
	}
}

func assertTestSignIn(s *SuiteController, testCase requestTestCase) {
	// config result
	want := testCase.responseJSONs
	// mock database
	if testCase.expected == true {
		s.mock.ExpectQuery(
			"SELECT `username`,`password` FROM `users` WHERE `username` = ?").
			WithArgs(testCase.input.Username).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password"}).
				AddRow(testCase.databaseReturn[0].Username, testCase.databaseReturn[0].Password))

	} else {
		s.mock.ExpectQuery(
			"SELECT `username`,`password` FROM `users` WHERE `username` = ?").
			WithArgs(testCase.input.Username).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password"}))
	}

	// mock a request
	r := gin.Default()
	r.POST("auth/sign-in", s.server.SignIn)
	w := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(testCase.input)
	req, err := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(jsonValue))
	assert.Nil(s.T(), err)
	r.ServeHTTP(w, req)
	assert.Equal(s.T(), testCase.statusCode, w.Code)

	var got map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), want, got)

}

var signUpTestCases = []requestTestCase{
	{
		testName: "success to sign up",
		expected: true,
		input:    models.User{FullName: "admin1234" ,Username: "admin1234", Password: "81dc9bdb52d04dc20036dbd8313ed055", AuthorityType: "customer"},
		databaseReturn: nil,
		responseJSONs: map[string]interface{}{
			"success": true,
			"message": "Sign up success",
		},
		statusCode: http.StatusOK,
	},
	{
		testName: "username is already taken",
		expected: false,
		input:    models.User{FullName: "admin1234" ,Username: "admin1234", Password: "81dc9bdb52d04dc20036dbd8313ed055", AuthorityType: "customer"},
		databaseReturn: []models.User{
			{Username: "admin1234", Password: "81dc9bdb52d04dc20036dbd8313ed055"},
		},
		responseJSONs: map[string]interface{}{
			"success": false,
			"message": "Username is already taken",
		},
		statusCode: http.StatusBadRequest,
	},
}

func (s *SuiteController) TestSignUp() {
	for _, c := range signUpTestCases {
		s.T().Run(c.testName, func(t *testing.T) {
			assertTestSignUp(s, c)
		})
	}
}

func assertTestSignUp(s *SuiteController, testCase requestTestCase) {
	// config result
	want := testCase.responseJSONs
	// mock database
	if testCase.expected == true {
		s.mock.MatchExpectationsInOrder(false)
		s.mock.ExpectBegin()

		s.mock.ExpectQuery(
			"SELECT `username`,`password` FROM `users` WHERE `username` = ?").
			WithArgs(testCase.input.Username).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password"}))

		s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
			WithArgs(testCase.input.FullName, testCase.input.Username, testCase.input.Password, testCase.input.AuthorityType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		s.mock.ExpectCommit()

	} else {
		s.mock.ExpectQuery(
			"SELECT `username`,`password` FROM `users` WHERE `username` = ?").
			WithArgs(testCase.input.Username).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password"}).AddRow(testCase.databaseReturn[0].Username, testCase.databaseReturn[0].Password))
			
	}

	// mock a request
	r := gin.Default()
	r.POST("auth/sign-up", s.server.SignUp)
	w := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(testCase.input)
	req, err := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(jsonValue))
	assert.Nil(s.T(), err)
	r.ServeHTTP(w, req)
	assert.Equal(s.T(), testCase.statusCode, w.Code)

	var got map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	fmt.Println(got)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), want, got)

}
