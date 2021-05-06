package dbhandle

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/config"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/infra/waf/logger"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type DBHandle struct {
	DBInf *gorm.DB
}

func (dbh DBHandle) ConInf() *gorm.DB {
	return dbh.DBInf
}

func NewDBHandler(conf config.Config) *DBHandle {
	var DB *gorm.DB
	var err error

	switch os.Getenv("APP_MODE") {
	case "production", "staging":
		DB, err = connectHerokuDB()
	case "localhost":
		DB, err = connectLocalDB(conf)
	default:
		DB, err = connectLocalDB(conf)
	}

	if err != nil {
		logger.Log.FatalWithFields("DB接続中にエラーが発生しました", database.Fields{"message": err.Error()})
	}
	DB.AutoMigrate(&domain.User{}, &domain.Bosyu{}, &domain.Message{})
	return &DBHandle{DB}
}

func connectLocalDB(conf config.Config) (*gorm.DB, error) {
	DBUrl := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable", conf.DB.Host, conf.DB.User, conf.DB.Name, conf.DB.Password, conf.DB.Port)
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
