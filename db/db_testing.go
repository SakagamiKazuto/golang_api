package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"work/model"
)


func ConnectTestDB() {
	DBUser     := "root"
	DBPass     := "root"
	DBProtocol := "tcp(db:3306)"
	DBName     := "golang_mysql_test"
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

	ts.Create(
		&model.User{
			Name: "User1",
			Mail: "sample1@gmail.com",
		},
	)

	ts.Create(
		&model.Bosyu{
			Title:  "明示的にUSERIDを指定2",
			UserID: 1,
		},
	)

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
}