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

	b = model.CreateBosyu(b, db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestCreateBosyuModel: %v", err)
	}
}

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

	title := "sample1_title_updated"
	about := "sample1_about_updated"
	pref:= "sample1_pref_updated"
	city := "sample1_city_updated"
	level := "sample1_level_updated"
	user_id := 1
	id := 3
	bosyu_json := fmt.Sprintf("{\"title\": \"%v\", \"about\": \"%v\", \"pref\": \"%v\", \"city\": \"%v\", \"level\": \"%v\", \"user_id\": %v, \"id\": %v}", title, about, pref, city, level, user_id, id)
	req, rec := createUpdateRequest(bosyu_json, token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.UpdateBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUpdateBosyuHandlerError(t *testing.T) {
	e := echo.New()
	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}
	// 1.DBにBosyuのIDが存在しない
	id := 0
	title := "sample1_title_updated"
	about := "sample1_about_updated"
	pref:= "sample1_pref_updated"
	city := "sample1_city_updated"
	level := "sample1_level_updated"
	user_id := 1
	bosyu_json := fmt.Sprintf("{\"title\": \"%v\", \"about\": \"%v\", \"pref\": \"%v\", \"city\": \"%v\", \"level\": \"%v\", \"user_id\": %v, \"id\": %v}", title, about, pref, city, level, user_id, id)
	req, rec := createUpdateRequest(bosyu_json, token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.UpdateBosyu)(contents)
	res := exec

	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}

	// 2.TitleやAboutが空欄
	title = ""
	about = ""
	id = 1
	bosyu_json = fmt.Sprintf("{\"title\": \"%v\", \"about\": \"%v\", \"pref\": \"%v\", \"city\": \"%v\", \"level\": \"%v\", \"user_id\": %v, \"id\": %v}", title, about, pref, city, level, user_id, id)
	req, rec = createUpdateRequest(bosyu_json, token)

	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.UpdateBosyu)(contents)
	res = exec

	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusBadRequest, code)
	}

	//3. JWTの認証が未通過
	token, err = handler.CreateToken(uint(9999),"DONTEXIST@gmail.com")
	req, rec = createPostRequest(bosyu_json, token)
	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.UpdateBosyu)(contents)
	res = exec
	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}
}
//

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

	req, rec := createDeleteRequest("1", token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.DeleteBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestDeleteBosyuHandlerError(t *testing.T) {
	e := echo.New()
	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}
	// 1.DBにBosyuのIDが存在しない
	req, rec := createDeleteRequest("0", token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.DeleteBosyu)(contents)
	res := exec

	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}

	//2. bosyu_idが正しくない値
	req, rec = createDeleteRequest("invalid_param", token)

	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.DeleteBosyu)(contents)
	res = exec

	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusBadRequest, code)
	}

	//3. JWTの認証が未通過
	token, err = handler.CreateToken(uint(9999),"DONTEXIST@gmail.com")
	req, rec = createPostRequest("1", token)
	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.DeleteBosyu)(contents)
	res = exec
	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusNotFound, code)
	}
}

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

func createUpdateRequest(bosyu_json string, token string) (*http.Request, *httptest.ResponseRecorder) {
	bodyReader := strings.NewReader(bosyu_json)
	req := httptest.NewRequest("PUT", "/api/bosyu/update", bodyReader)
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

func createDeleteRequest(bID string, token string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/bosyu/delete?bosyu_id=%v", bID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	return req, rec
}


func getErrorStatusCode(res interface{}) int {
	code := reflect.Indirect(reflect.ValueOf(res)).Field(0).Interface().(int)
	return code
}
