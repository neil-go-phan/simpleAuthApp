package models

import (
	"database/sql"
	"regexp"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository Repository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err)

	conn, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	s.DB = conn
	require.NoError(s.T(), err)

	s.repository = CreateRepository(s.DB)
}
func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}
func (s *Suite) Test_repository_GetUsers() {
	rows := sqlmock.NewRows([]string{"full_name", "username", "password", "authority_type"}).
		AddRow("admin", "admin1234", "81dc9bdb52d04dc20036dbd8313ed055", "admin").
		AddRow("customer", "customer1234", "81dc9bdb52d04dc20036dbd8313ed055", "customer").
		AddRow("employee", "employee1234", "81dc9bdb52d04dc20036dbd8313ed055", "employee")
	// TODO: search more about this f*cking ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnRows(rows)

	_, err := s.repository.GetUsers()
	assert.Nil(s.T(), err)

}

func (s *Suite) Test_repository_FindUser() {
	want := User{
		Username: "admin1234",
		Password: "81dc9bdb52d04dc20036dbd8313ed055",
	}
	s.mock.ExpectQuery(
		"SELECT `username`,`password` FROM `users` WHERE `username` = ?").
		WithArgs(want.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "password"}).
			AddRow(want.Username, want.Password))

	got, err := s.repository.FindUser(want.Username)
	assert.Nil(s.T(), err)

	if *got != want {
		s.T().Error("Query result is wrong")
	}

}

func (s *Suite) Test_repository_CreateUser() {
	want := User{
		FullName:      "test username",
		Username:      "testuser",
		Password:      "81dc9bdb52d04dc20036dbd8313ed055",
		AuthorityType: "customer",
	}
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(want.FullName, want.Username, want.Password, want.AuthorityType).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	_, err := s.repository.CreateUser(&want)
	assert.Nil(s.T(), err)

}
