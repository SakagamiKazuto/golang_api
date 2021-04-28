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


func NewDBHandler() *DBHandle {
	var DB *gorm.DB
	var err error

	if os.Getenv("DATABASE_URL") != "" {
		DB, err = connectHerokuDB()
	} else {
		DB, err = connectLocalDB()
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	DB.AutoMigrate(&domain.User{}, &domain.Bosyu{}, &domain.Message{})
	return &DBHandle{DB}
}

func connectLocalDB() (*gorm.DB, error) {
	err := godotenv.Load("/go/src/.env")
	if err != nil {
		return nil, err
	}

	DBHost := os.Getenv("DB_HOST")
	DBUser := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_NAME")
	DBPass := os.Getenv("DB_PASSWORD")
	DBPort := os.Getenv("DB_PORT")
	DBUrl := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", DBHost, DBUser, DBName, DBPass, DBPort)
	return connectDB(DBUrl)
}

func connectHerokuDB() (*gorm.DB, error) {
	DBUrl := os.Getenv("DATABASE_URL")
	return connectDB(DBUrl)
}

func connectDB(DBUrl string) (*gorm.DB, error) {
	d := "postgres"
	db, err := gorm.Open(d, DBUrl)
	return db, err
}
