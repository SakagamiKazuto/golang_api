package test

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewTestDB() *gorm.DB {
	if err := godotenv.Load("/go/src/.env"); err != nil {
		throughError(err)
	}

	DBHost := os.Getenv("DB_HOST")
	DBUser := os.Getenv("DB_USER")
	DBName := os.Getenv("TEST_DB_NAME")
	DBPass := os.Getenv("DB_PASSWORD")
	DBPort := os.Getenv("DB_PORT")
	CONNECT := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", DBHost, DBUser, DBName, DBPass, DBPort)
	db, err := gorm.Open("postgres", CONNECT)
	if err != nil {
		throughError(err)
	}
	db.AutoMigrate(&domain.User{}, &domain.Bosyu{}, &domain.Message{})
	return db
}

func InsertTestData(db *gorm.DB) {
	ts := db.Begin()
	defer ts.Commit()

	// User Data
	ts.Create(&domain.User{Name: "sample1", Mail: "sample1@gmail.com", Password: "123", Model: gorm.Model{ID: 1}})

	// Bosyu Data
	// normal pattern1
	ts.Create(&domain.Bosyu{Title: "sample1", About: "sample1", UserID: 1, Model: gorm.Model{ID: 1}})
	// deleted_at is not Null
	ts.Create(&domain.Bosyu{Title: "sample2", About: "sample2", UserID: 1, Model: gorm.Model{ID: 2, DeletedAt: getTimeNowPointer()}})
	// normal pattern2
	ts.Create(&domain.Bosyu{Title: "sample3", About: "sample3", UserID: 1, Model: gorm.Model{ID: 3}})

	if err := ts.Error; err != nil {
		ts.Rollback()
	}
	return
}

func DeleteTestData(db *gorm.DB) {
	ts := db.Begin()
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

func throughError(err error) string {
	panic(fmt.Sprintf("エラーが発生しました\n%s", err.Error()))
}
