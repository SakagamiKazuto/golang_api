package test

import (
	"fmt"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"

	"work/common"
	"work/db"
	"work/handler"
	"work/model"
)


/*
CreateBosyuTests
Model:
1. 正しいパターンではSQLのクエリが正しく発行されているかのチェック
Handler:
Normal
1. status201

Error
1. TitleとAboutの空欄はエラー
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

	b = model.CreateBosyu(b, db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestCreateBosyuModel: %v", err)
	}
}

// 1. 正しい値のときはJSONが帰ってくる
func TestCreateBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	bosyu_json := `{"title": "sample_title", "about": "sample_about", "pref": "愛媛県", "city": "松山市", "level": "player", "user_id": 123123}`
	token, err := createTokenFromSomeUser()
		if err != nil {
			t.Errorf("got error like: %+v", err)
		}

	req, rec := createPostRequest(bosyu_json, token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)

	assert.NoError(t, exec)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

//2. Title or Aboutにおける空欄ではエラーを返す
func TestCreateBosyuHandlerError(t *testing.T) {
	e := echo.New()
	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	// Titleが空欄
	bosyu_json := `{"title": "", "about": "sample_about", "pref": "愛媛県", "city": "松山市", "level": "player", "user_id": 123123}`
	req, rec := createPostRequest(bosyu_json, token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)
	res := exec

	if assert.Error(t, res) {
		//	res.Codeだと値が取り出せないので
		//	handlerがreturnしてきた値からstatusCodeを取り出す
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusBadRequest, code)
	}

	// Aboutが空欄
	bosyu_json = `{"title": "sample_title", "about": "", "pref": "愛媛県", "city": "松山市", "level": "player", "user_id": 123123}`
	req, rec = createPostRequest(bosyu_json, token)

	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)
	res = exec

	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusBadRequest, code)
	}

	//3. JWTの認証が未通過
	token, err = handler.CreateToken(uint(9999),"DONTEXIST@gmail.com")
	req, rec = createPostRequest(bosyu_json, token)
	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)
	res = exec
	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}
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
	userId := uint(1)
	bosyus := model.FindBosyu(userId, db.DB)
	for _, bosyu := range bosyus {
		assert.Empty(t, nil, bosyu.DeletedAt)
		assert.Equal(t, userId, bosyu.UserID)
	}
}

func TestGetBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	req, rec := createGetRequest("1", token)
	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.GetBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}


func TestGetBosyuHandlerError(t *testing.T) {
	e := echo.New()

	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	//1. データがDBに存在しない
	req, rec := createGetRequest("99999", token)
	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.GetBosyu)(contents)
	res := exec
	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}


	//2. user_idの値がBlank
	req, rec = createGetRequest("", token)
	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.GetBosyu)(contents)
	res = exec
	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusBadRequest, code)
	}

	//3. JWTの認証が未通過
	token, err = handler.CreateToken(uint(9999),"DONTEXIST@gmail.com")
	req, rec = createGetRequest("1", token)
	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.GetBosyu)(contents)
	res = exec
	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}
}

func TestUpdateBosyu(t *testing.T) {
	e := echo.New()

	// テスト用リクエスト生成
	bosyu_json := `{"id": 1, "title": "sample_title", "about": "sample_about", "pref": "香川県", "city": "高松市", "level": "player", "user_id": 1}`
	bodyReader := strings.NewReader(bosyu_json)
	req := httptest.NewRequest("PUT", "/api/bosyu/update", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()

	contents := e.NewContext(req, rec)

	assert.NoError(t, handler.UpdateBosyu(contents))
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "高松市")
}
//
//// !! 1.存在するID指定 2. IDが存在しないパターン 3. 0になったときerror吐くパターン
//func TestDeleteBosyu(t *testing.T) {
//	// echo.New()をテストの中で何度も呼ばなくて住むようにしたい
//	e := echo.New()
//
//	// user_id exists in Record pattern
//	req := httptest.NewRequest("DELETE", "/api/bosyu/delete?id=1", nil)
//	rec := httptest.NewRecorder()
//
//	contents := e.NewContext(req, rec)
//
//	assert.NoError(t, handler.DeleteBosyu(contents))
//	assert.Equal(t, http.StatusNoContent, rec.Code)
//}

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

func createTokenFromSomeUser()(string, error) {
	user := model.FindUser(&model.User{}, db.DB)
	token, err := handler.CreateToken(user.ID, user.Mail)
	return token, err
}

func createPostRequest(bosyu_json string, token string) (*http.Request, *httptest.ResponseRecorder) {
	bodyReader := strings.NewReader(bosyu_json)
	req := httptest.NewRequest("POST", "/api/bosyu/create", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	return req, rec
}

func createGetRequest(uID string, token string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/bosyu/get?user_id=%v", uID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	return req, rec
}


func getErrorStatusCode(res interface{}) int {
	code := reflect.Indirect(reflect.ValueOf(res)).Field(0).Interface().(int)
	return code
}
