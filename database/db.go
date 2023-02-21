package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DB_USERNAME = "root"
const DB_PASSWORD = "12321221"
const DB_NAME = "store"
const DB_HOST = "localhost"
const DB_PORT = "3306"

type Database interface {

}

type GROMDatabase struct {
	*gorm.DB
}

var Db *gorm.DB


func ConnectDB() {
	var err error
	dsn := DB_USERNAME +":"+ DB_PASSWORD +"@tcp"+ "(" + DB_HOST + ":" + DB_PORT +")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
	}
	Db = db
}

func GetDB() *gorm.DB {
	return Db
}