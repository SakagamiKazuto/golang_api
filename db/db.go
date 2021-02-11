package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"os"
	"work/model"
)

var DB *gorm.DB

const Dialect = "mysql"

func InitDB() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		panic("failed to read .env")
	}

	DB, err = connectDB()
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&model.User{}, &model.Bosyu{}, &model.Message{})
}

func connectDB() (*gorm.DB, error) {
	var DBUser string
	var DBPass string
	var DBProtocol string
	var DBName string
	if os.Getenv("CLEARDB_DATABASE_URL") == "" {
		DBUser     = os.Getenv("LOCAL_USER")
		DBPass     = os.Getenv("LOCAL_PASSWORD")
		DBProtocol = os.Getenv("LOCAL_PROTOCOL")
		DBName     = os.Getenv("LOCAL_DBNAME")
	} else {
		DBUser     = os.Getenv("HEROKU_USER")
		DBPass     = os.Getenv("HEROKU_PASSWORD")
		DBProtocol = os.Getenv("HEROKU_PROTOCOL")
		DBName     = os.Getenv("HEROKU_DBNAME")
	}
	CONNECT := DBUser + ":" + DBPass + "@" + DBProtocol + "/" + DBName + "?parseTime=true"
	db, err := gorm.Open(Dialect, CONNECT)
	return db, err
}

