package test

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/common"
	"github.com/SakagamiKazuto/golang_api/db"
	"github.com/SakagamiKazuto/golang_api/handler"
	"github.com/SakagamiKazuto/golang_api/model"
)

/*
CreateBosyuTests
Model:
1. 正しいパターンではSQLのクエリが正しく発行されているかのチェック
Handler:
Normal
1. status201

Error
1. TitleとAboutの空欄
2. JWTの認証が通らない
*/
func TestCreateBosyuModel(t *testing.T) {
	db, mock, err := getDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
		return
	}

	b := new(model.Bosyu)

	b.Title = "sample_title"
	b.About = "sample_about"
	b.Prefecture = "愛媛県"
	b.City = "松山市"
	b.Level = "player"
	b.UserID = 123123

	// 発行されているクエリおよび影響を受けるカラムの数が正しいかをチェックする
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `bosyus` (`created_at`,`updated_at`,`deleted_at`,`title`,`about`,`prefecture`,`city`,`level`,`user_id`) VALUES (?,?,?,?,?,?,?,?,?)")).
		WithArgs(common.AnyTime{}, common.AnyTime{}, nil, b.Title, b.About, b.Prefecture, b.City, b.Level, b.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	b, err = model.CreateBosyu(b, db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestCreateBosyuModel: %v", err)
	}
}

func TestCreateBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	title := "sample_title"
	about := "sample_about"
	pref := "sample_pref"
	city := "sample_city"
	level := "sample_level"
	user_id := 1
	bosyuJson := fmt.Sprintf("{\"title\": \"%v\", \"about\": \"%v\", \"pref\": \"%v\", \"city\": \"%v\", \"level\": \"%v\", \"user_id\": %v}", title, about, pref, city, level, user_id)
	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	mockReq := MockReq{bosyuJson, token, "/api/bosyu/create",  "POST"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestCreateBosyuHandlerError(t *testing.T) {
}

/*
GetBosyuTests
Model:
1. テストデータにdeleted_atがnullとそうでないものを用意し、deleted_atがnullのものだけ返すことを確認
Handler:
Normal
1. status200
Error
1. データがDBに存在しない
2. user_idの値がBlank
3. JWT認証が通らない
*/
func TestGetBosyuModel(t *testing.T) {
}

func TestGetBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	mockReq := MockReq{`{"user_id": 1}`, token, "/api/bosyu/get", "GET"}
	req, rec := mockReq.createReq()
	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.GetBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetBosyuHandlerError(t *testing.T) {
}

/*
UpdateBosyuTests
Model:
Normal
1. BosyuのIDがdatabaseに存在する

Error
1. BosyuのIDがdatabaseに存在しない

Handler:
Normal
1. status200

Error
1. DBにIDが一致する募集が存在しない
2. TitleとAboutの空欄である
3. JWTの認証が通らない
*/
func TestUpdateBosyuModelNormal(t *testing.T) {
	b := new(model.Bosyu)

	b.Title = "sample1_title_updated"
	b.About = "sample1_about_updated"
	b.Prefecture = "sample1_pref_updated"
	b.City = "sample1_city_updated"
	b.Level = "sample1_level_updated"
	b.UserID = 1
	b.ID = 1
	bosyu, err := model.UpdateBosyu(b, db.DB)
	if assert.NoError(t, err) {
		assert.Empty(t, nil, bosyu.DeletedAt)
		assert.Equal(t, "sample1_title_updated", bosyu.Title)
	}
}

func TestUpdateBosyuModelError(t *testing.T) {
	b := new(model.Bosyu)

	b.Title = "sample1_title_updated"
	b.About = "sample1_about_updated"
	b.Prefecture = "sample1_pref_updated"
	b.City = "sample1_city_updated"
	b.Level = "sample1_level_updated"
	b.UserID = 1
	b.ID = 0

	_, err := model.UpdateBosyu(b, db.DB)
	assert.Error(t, err)

}

func TestUpdateBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	title := "sample3_title_updated"
	about := "sample3_about_updated"
	pref := "sample3_pref_updated"
	city := "sample3_city_updated"
	level := "sample3_level_updated"
	user_id := 1
	id := 3
	bosyuJson := fmt.Sprintf("{\"title\": \"%v\", \"about\": \"%v\", \"pref\": \"%v\", \"city\": \"%v\", \"level\": \"%v\", \"user_id\": %v, \"id\": %v}", title, about, pref, city, level, user_id, id)
	mockReq := MockReq{bosyuJson, token, "/api/bosyu/update",  "PUT"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.UpdateBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUpdateBosyuHandlerError(t *testing.T) {
}

/*
DeleteBosyuTests
Model:
Normal
1. BosyuのIDがdatabaseに存在する

Error
1. BosyuのIDがdatabaseに存在しない

Handler:
Normal
1. status200

Error
1. DBにIDが一致する募集が存在しない
2. bosyu_idが正しくない値
3. JWTの認証が通らない
*/
func TestDeleteBosyuModelNormal(t *testing.T) {
	b := new(model.Bosyu)
	b.ID = 1
	err := model.DeleteBosyu(b.ID, db.DB)
	assert.NoError(t, err)
}

func TestDeleteBosyuModelError(t *testing.T) {
	b := new(model.Bosyu)
	b.ID = 0

	err := model.DeleteBosyu(b.ID, db.DB)
	assert.Error(t, err)
}

func TestDeleteBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	mockReq := MockReq{`{"id": 1}`, token, "/api/bosyu/delete",  "DELETE"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.DeleteBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteBosyuHandlerError(t *testing.T) {
}

// CommonMethod's
func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

func createTokenFromSomeUser() (string, error) {
	user, err := model.FindUser(&model.User{}, db.DB)
	token, err := handler.CreateToken(user.ID, user.Mail)
	return token, err
}

func getErrorStatusCode(res interface{}) int {
	code := reflect.Indirect(reflect.ValueOf(res)).Field(0).Interface().(int)
	return code
}
