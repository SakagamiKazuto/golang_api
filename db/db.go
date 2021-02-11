package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"os"
	"work/model"
)

var DB *gorm.DB
var err error

const Dialect = "mysql"

func init() {
	err = godotenv.Load(".env")
	if err != nil {
		panic("failed to read .env")
	}

	DB, err = connectDB(err)
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&model.User{}, &model.Bosyu{}, &model.Message{})
}

func connectHerokuDB(err error) (*gorm.DB, error) {
	CONNECT := "b4ffeaf9d21ec3:e23d209c@tcp(us-cdbr-east-03.cleardb.com)/heroku_3afafbe663f4456?parseTime=true"
	DB, err = gorm.Open(Dialect, CONNECT)
	return DB, err
}

func connectDB(err error) (*gorm.DB, error) {
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
	DB, err = gorm.Open(Dialect, CONNECT)
	return DB, err
}
//func insertSmapleData(db *gorm.DB) {
//	ts := db.Begin()
//	defer ts.Commit()
//
//	ts.Create(
//		&User{
//			Name: "SosikiA",
//			Bosyu: []Bosyu{
//				{Title: "明示的にUSERIDを指定", UserID: 1},
//			},
//		},
//	)
//
//	ts.Create(
//		&Bosyu{
//			Title:  "明示的にUSERIDを指定2",
//			UserID: 1,
//		},
//	)
//
//	if err := ts.Error; err != nil {
//		ts.Rollback()
//	}
//	return
//}
