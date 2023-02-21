package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectingDatabase(t *testing.T) {
	assert := assert.New(t)
	ConnectDB()
	db := GetDB()
	sqlDB, _:= db.DB()
	assert.NoError(sqlDB.Ping(), "DB should be able to ping")
	sqlDB.Close()
}

