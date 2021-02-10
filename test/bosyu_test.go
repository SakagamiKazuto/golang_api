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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"

	"work/common"
	"work/model"
	"work/handler"
	//"work/db"
)

/* Tasks
!!DB初期化を保証できるバッチ関数を作ることでDBに入ってる値を保証したい
1.recordの削除 2.DB値の初期化を処理として入れたい → DB.createとか使えばいけそう

!!各API例外パターンテスト & 実装したい
!!結果の検証がassert.containsで検証しているだけで、このままだと"余計なもの含んでるパターン"検証できないので細分化したい
!!handle, modelでテストを分割したい
*/

/*
CreateBosyuTests
Model:
1. 正しいパターンではSQLのクエリが正しく発行されているかのチェック
Handler:
1. 正しい値のときはJSONが帰ってくる
2. TitleとAboutの空欄はエラー
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
func TestCreateBosyuHandlerNormalPattern(t *testing.T) {
	e := echo.New()

	bosyu_json := `{"title": "sample_title", "about": "sample_about", "pref": "愛媛県", "city": "松山市", "level": "player", "user_id": 123123}`
	token, err := handler.CreateToken(uint(3), "sample2@mail.com")
		if err != nil {
			t.Errorf("got error like: %+v", err)
		}

	req, rec := CreatePostRequest(bosyu_json, token)

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)

	assert.NoError(t, exec)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

//2. Title or Aboutにおける空欄ではエラーを返す
func TestCreateBosyuHandlerErrorPattern(t *testing.T) {
	e := echo.New()
	token, err := handler.CreateToken(uint(3), "sample2@mail.com")
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	// Titleが空欄
	bosyu_json := `{"title": "", "about": "sample_about", "pref": "愛媛県", "city": "松山市", "level": "player", "user_id": 123123}`
	req, rec := CreatePostRequest(bosyu_json, token)

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
	req, rec = CreatePostRequest(bosyu_json, token)

	contents = e.NewContext(req, rec)
	exec = middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)
	res = exec

	if assert.Error(t, res) {
		code := getErrorStatusCode(res)
		assert.Equal(t, http.StatusBadRequest, code)
	}
}

func getErrorStatusCode(res interface{}) int {
	code := reflect.Indirect(reflect.ValueOf(res)).Field(0).Interface().(int)
	return code
}

//// !! 1.id指定 4.それぞれblank時点の動作をテストしたい
//// 2.住所指定 3.市指定はapiの仕様上保留
//func TestGetBosyu(t *testing.T) {
//	// echo.New()をテストの中で何度も呼ばなくて住むようにしたい
//	e := echo.New()
//
//	// user_id exists in Record pattern
//	req := httptest.NewRequest("GET", "/api/bosyu/get?user_id=123123", nil)
//	rec := httptest.NewRecorder()
//
//	contents := e.NewContext(req, rec)
//
//	assert.NoError(t, handler.GetBosyu(contents))
//	assert.Equal(t, http.StatusOK, rec.Code)
//	assert.Contains(t, rec.Body.String(), "\"user_id\":123123")
//
//	// 2. user_id is blank pattern
//	req = httptest.NewRequest("GET", "/api/bosyu/get", nil)
//	rec = httptest.NewRecorder()
//
//	contents = e.NewContext(req, rec)
//
//	assert.NoError(t, handler.GetBosyu(contents))
//	assert.Equal(t, http.StatusOK, rec.Code)
//	assert.Contains(t, rec.Body.String(), "\"user_id\":1")
//
//
//	// 3. user_id doesnot exist in Record pattern
//	req = httptest.NewRequest("GET", "/api/bosyu/get?user_id=2", nil)
//	rec = httptest.NewRecorder()
//
//	contents = e.NewContext(req, rec)
//
//	assert.NoError(t, handler.GetBosyu(contents))
//	assert.Equal(t, http.StatusOK, rec.Code)
//	assert.Contains(t, rec.Body.String(), "[]")
//}
//
//func TestUpdateBosyu(t *testing.T) {
//	e := echo.New()
//
//	// テスト用リクエスト生成
//	bosyu_json := `{"id": 1, "title": "sample_title", "about": "sample_about", "pref": "香川県", "city": "高松市", "level": "player", "user_id": 1}`
//	bodyReader := strings.NewReader(bosyu_json)
//	req := httptest.NewRequest("PUT", "/api/bosyu/update", bodyReader)
//	req.Header.Add("Content-Type", "application/json")
//	req.Header.Add("Accept", "application/json")
//	rec := httptest.NewRecorder()
//
//	contents := e.NewContext(req, rec)
//
//	assert.NoError(t, handler.UpdateBosyu(contents))
//	assert.Equal(t, http.StatusCreated, rec.Code)
//	assert.Contains(t, rec.Body.String(), "高松市")
//}
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

func CreatePostRequest(bosyu_json string, token string) (*http.Request, *httptest.ResponseRecorder) {
	bodyReader := strings.NewReader(bosyu_json)
	req := httptest.NewRequest("POST", "/api/bosyu/create", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	return req, rec
}

