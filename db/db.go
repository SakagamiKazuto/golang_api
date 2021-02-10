package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"work/model"
)

const (
	Dialect    = "mysql"
	DBUser     = "root"
	DBPass     = "root"
	DBProtocol = "tcp(db:3306)"
	DBName     = "matching_portfolio"
)

var DB *gorm.DB

func init() {
	var err error
	//CONNECT := DBUser + ":" + DBPass + "@" + DBProtocol + "/" + DBName + "?parseTime=true"
	//DB, err = gorm.Open(Dialect, CONNECT)
	DB, err = connectHerokuDB(err)
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&model.User{}, &model.Bosyu{}, &model.Message{})
}

func connectHerokuDB(err error) (*gorm.DB, error) {
	CONNECT := "b4ffeaf9d21ec3:e23d209c@us-cdbr-east-03.cleardb.com/heroku_3afafbe663f4456?parseTime=true"
	DB, err = gorm.Open(Dialect, CONNECT)
	return DB, err
}

func connectLocalDB(err error) (*gorm.DB, error) {
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
