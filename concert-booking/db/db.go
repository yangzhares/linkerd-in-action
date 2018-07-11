package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB struct {
	*gorm.DB
}

func InitDB(dbname, user, password, endpoint string) (*DB, error) {
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, endpoint, dbname)
	db, err := gorm.Open("mysql", dbinfo)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
		return nil, err
	}

	return &DB{DB: db}, nil
}
