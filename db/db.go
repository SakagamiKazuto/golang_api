package db

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var DB *gorm.DB

const Dialect = "mysql"

func InitDB() {

	DB, err := connectDB()
	if err != nil {
		//panic("failed to connect database")
		panic(err.Error())
	}
	DB.AutoMigrate(&model.User{}, &model.Bosyu{}, &model.Message{})
}

func connectDB() (*gorm.DB, error) {
	var CONNECT string
	if os.Getenv("DATABASE_URL") != "" {
		CONNECT = os.Getenv("DATABASE_URL")
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			panic("failed to read .env")
		}

		DBUser := os.Getenv("LOCAL_USER")
		DBPass := os.Getenv("LOCAL_PASSWORD")
		DBName := os.Getenv("LOCAL_DBNAME")
		CONNECT = fmt.Sprintf("host=db user=%s dbname=%s password=%s port=5432 sslmode=disable", DBUser, DBName, DBPass)
	}
	db, err := gorm.Open("postgres", CONNECT)

	return db, err
}
