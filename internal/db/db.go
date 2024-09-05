package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {
	dburl := os.Getenv("DB_URL")
	var err error

	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		panic("failed to connect to database")
	}

	query := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""
	if err = DBConn.Exec(query).Error; err != nil {
		fmt.Println("cannot install uuid extension")
		panic(err)
	}

	// Migrate Models
	if err = DBConn.AutoMigrate(&User{}, &SearchSettings{}, &CrawledUrl{}); err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DBConn
}    