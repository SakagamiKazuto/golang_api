package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/handler"
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
	req, rec := createUserSignupRequest(user_json)

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
	req, rec := createUserSignupRequest(user_json)

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
	req, rec = createUserSignupRequest(user_json)

	contents = e.NewContext(req, rec)
	res = handler.Signup(contents)
	if assert.Error(t, res) {
		assert.Equal(t, http.StatusConflict, getErrorStatusCode(res))
	}
}


/*
LoginTests
Handler:
Normal
1. status200

Error
1. Mail, Passwordの値に基づくUserが存在しない
*/
func TestLoginNormal(t *testing.T) {
	e := echo.New()

	Mail := "sample1@gmail.com"
	Password := "123"
	Name := "sample_name"
	Address:= "sample_address"
	Tel := "sample_tel"
	user_json := fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, Name, Address, Tel, Mail)
	req, rec := createUserLoginRequest(user_json)

	contents := e.NewContext(req, rec)
	if assert.NoError(t, handler.Login(contents)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestLoginError(t *testing.T) {
	e := echo.New()

	// 1.Mail, Passwordの値に基づくUserが存在しない
	Mail := "not_exist@gmail.com"
	Password := "9999999999"
	Name := "sample_name"
	Address:= "sample_address"
	Tel := "sample_tel"
	user_json := fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, Name, Address, Tel, Mail)
	req, rec := createUserLoginRequest(user_json)

	contents := e.NewContext(req, rec)
	if assert.Error(t, handler.Login(contents)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}


// CommonMethod's
func createUserSignupRequest(user_json string) (*http.Request, *httptest.ResponseRecorder) {
	bodyReader := strings.NewReader(user_json)
	req := httptest.NewRequest("POST", "/signup", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	return req, rec
}

func createUserLoginRequest(user_json string) (*http.Request, *httptest.ResponseRecorder) {
	bodyReader := strings.NewReader(user_json)
	req := httptest.NewRequest("POST", "/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	return req, rec
}




