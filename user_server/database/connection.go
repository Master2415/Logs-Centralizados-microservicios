package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var StartTime time.Time
var DB *gorm.DB

func Connection() {
	var host, user, password, dbname, port string

	host = os.Getenv("US_DATABASE_URL")
	user = os.Getenv("US_USER")
	password = os.Getenv("US_PASSWORD")
	dbname = os.Getenv("US_DBNAME")
	port = os.Getenv("US_PORT")

	if host == "" {
		host = "localhost"
		user = "admin"
		password = "12345"
		dbname = "users_db"
		port = "5432"
	}

	var DSN = "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port

	var error error
	StartTime = time.Now()
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Fatal(error)
	} else {
		log.Println("USERs Database => " + DSN + "\n\n")
	}
}

func CheckDatabaseLive() bool {
	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		return false
	}

	return true
}

func CheckDatabaseReady() bool {
	sqlDB, err := DB.DB()
	err = sqlDB.Ping()
	if err != nil {
		return false
	}
	return true
}
