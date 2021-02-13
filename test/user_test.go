package test

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"work/handler"
)

/*
SignupTests
Handler:
Normal
1. status201

Error
1. Name, Mail, Passwordのいずれかが空欄
2. MailがすでにUsersのテーブルに存在する
*/
func TestSignupNormal(t *testing.T) {
	e := echo.New()

	Password := "sample_password"
	Name := "sample_name"
	Address:= "sample_address"
	Tel := "sample_tel"
	Mail := "sample_mail"
	user_json := fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, Name, Address, Tel, Mail)
	req, rec := createUserPostRequest(user_json)

	contents := e.NewContext(req, rec)
	if assert.NoError(t, handler.Signup(contents)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestSignupError(t *testing.T) {
	e := echo.New()

	// 1.Name, Mail, Passwordのいずれかが空欄
	Password := ""
	Name := ""
	Address:= "sample_address"
	Tel := "sample_tel"
	Mail := ""
	user_json := fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, Name, Address, Tel, Mail)
	req, rec := createUserPostRequest(user_json)

	contents := e.NewContext(req, rec)
	res := handler.Signup(contents)
	if assert.Error(t, res) {
		assert.Equal(t, http.StatusBadRequest, getErrorStatusCode(res))
	}

	//2.MailがすでにUsersのテーブルに存在する
	Mail = "sample1@gmail.com"
	Password = "sample_password"
	Name = "sample_name"
	user_json = fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, Name, Address, Tel, Mail)
	req, rec = createUserPostRequest(user_json)

	contents = e.NewContext(req, rec)
	res = handler.Signup(contents)
	if assert.Error(t, res) {
		assert.Equal(t, http.StatusConflict, getErrorStatusCode(res))
	}
}


//func TestLogin(t *testing.T) {
//	e := echo.New()
//
//	// ↑のテストでDBに追加済み
//	user_json := `{"Password": "dummy", "Name": "dummy_account", "address": "京都市","Tel": "000-0000-0000", "Mail": "sample@gmail.com"}`
//	bodyReader := strings.NewReader(user_json)
//	req := httptest.NewRequest("POST", "/login", bodyReader)
//	req.Header.Add("Content-Type", "application/json")
//	req.Header.Add("Accept", "application/json")
//	rec := httptest.NewRecorder()
//
//	contents := e.NewContext(req, rec)
//
//	assert.NoError(t, handler.Login(contents))
//	assert.Equal(t, http.StatusOK, rec.Code)
//	// 本当はJSONの中身にuser_jsonの値それぞれが含むかをチェックしたいが
//	// rec.Bodyに\などの記号が挿入されるのでこの程度のチェックで留める
//	assert.Contains(t, rec.Body.String(), "token")
//}


/*
	common methods
 */
func createUserPostRequest(user_json string) (*http.Request, *httptest.ResponseRecorder) {
	bodyReader := strings.NewReader(user_json)
	req := httptest.NewRequest("POST", "/signup", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	return req, rec
}





