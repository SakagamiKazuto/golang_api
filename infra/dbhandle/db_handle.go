package dbhandle

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type DBHandle struct {
	DBInf *gorm.DB
}

func (dbh DBHandle) ConInf() *gorm.DB {
	return dbh.DBInf
}

const Dialect = "postgres"

func NewDBHandler() *DBHandle {
	DB, err := connectDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	DB.AutoMigrate(&domain.User{}, &domain.Bosyu{}, &domain.Message{})
	return &DBHandle{DB}
}

func connectDB() (*gorm.DB, error) {
	var CONNECT string
	if os.Getenv("DATABASE_URL") != "" {
		CONNECT = os.Getenv("DATABASE_URL")
	} else {
		err := godotenv.Load("/go/src/.env")
		if err != nil {
			return nil, err
		}

		DBHost := os.Getenv("DB_HOST")
		DBUser := os.Getenv("DB_USER")
		DBName := os.Getenv("DB_NAME")
		DBPass := os.Getenv("DB_PASSWORD")
		DBPort := os.Getenv("DB_PORT")
		CONNECT = fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", DBHost, DBUser, DBName, DBPass, DBPort)
	}
	db, err := gorm.Open(Dialect, CONNECT)

	return db, err
}
