package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"work/model"
)

/*
自動テスト実行するにあたり必要な関数はこのファイルに置く
*/
func ConnectTestDB() {
	DBUser := "root"
	DBPass := "root"
	DBProtocol := "tcp(db:3306)"
	DBName := "golang_mysql_test"
	CONNECT := DBUser + ":" + DBPass + "@" + DBProtocol + "/" + DBName + "?parseTime=true"
	var err error
	DB, err = gorm.Open(Dialect, CONNECT)
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&model.User{}, &model.Bosyu{}, &model.Message{})
}

func InsertTestData() {
	ts := DB.Begin()
	defer ts.Commit()

	// User Data
	ts.Create(&model.User{Name: "sample1", Mail: "sample1@gmail.com", Password: "123", Model: gorm.Model{ID: 1}})

	// Bosyu Data
	// normal pattern1
	ts.Create(&model.Bosyu{Title: "sample1", About: "sample1", UserID: 1, Model: gorm.Model{ID:1}})
	// deleted_at is not Null
	ts.Create(&model.Bosyu{Title: "sample2", About: "sample2", UserID: 1, Model: gorm.Model{ID:2, DeletedAt: getTimeNowPointer()}})
	// normal pattern2
	ts.Create(&model.Bosyu{Title: "sample3", About: "sample3", UserID: 1, Model: gorm.Model{ID:3}})

	if err := ts.Error; err != nil {
		ts.Rollback()
	}
	return
}

func DeleteTestData() {
	ts := DB.Begin()
	defer ts.Commit()
	ts.Exec("DELETE FROM users")
	ts.Exec("DELETE FROM bosyus")
	if err := ts.Error; err != nil {
		ts.Rollback()
	}
	return
}

func getTimeNowPointer() *time.Time {
	now := time.Now()
	nowP := &now
	return nowP
}

