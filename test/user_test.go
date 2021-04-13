package test

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/db"
	"github.com/SakagamiKazuto/golang_api/model"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/handler"
)

/*
SignupTests
CreateUser:
Normal
1. データを作成する

Error
1. Mailの重複キー制約に違反
2. 内部エラー

Handler:
Normal
1. status201

Error
1. Name, Mail, Passwordのいずれかが空欄
2. MailがすでにUsersのテーブルに存在する
*/
func TestCreateUserNormal(t *testing.T) {
	u := new(model.User)
	u.Mail = "sample9@gmail.com"
	u.Password = "123123"
	_, err := model.CreateUser(u, db.DB)
	assert.NoError(t, err)
}

func TestCreateUserError(t *testing.T) {
	u := new(model.User)
	u.Mail = "sample1@gmail.com"
	_, err := model.CreateUser(u, db.DB)
	assert.Error(t, err, fmt.Sprintf(`メールアドレス%sのデータ挿入に失敗しました:pq: duplicate key value violates unique constraint "users_mail_key"`, u.Mail))
}

func TestSignupNormal(t *testing.T) {
	e := echo.New()

	Password := "sample_password"
	Name := "sample_name"
	Address := "sample_address"
	Tel := "sample_tel"
	Mail := "sample_mail"
	userJson := fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, Name, Address, Tel, Mail)
	mockReq := MockReq{userJson, "", "/signup", "POST"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	if assert.NoError(t, handler.Signup(contents)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestSignupError(t *testing.T) {
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
	userJson := fmt.Sprintf("{\"password\": \"%v\", \"name\": \"%v\", \"address\": \"%v\", \"tel\": \"%v\", \"mail\": \"%v\"}", Password, "", "", "", Mail)
	mockReq := MockReq{userJson, "", "/login", "POST"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	if assert.NoError(t, handler.Login(contents)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestLoginError(t *testing.T) {
}

/*
FindUserByUid
Normal
1. idに基づいてユーザー情報を取得

Error
1. 該当idのユーザーが存在しない
 */
func TestFindUserByUidNormal(t *testing.T) {
	u := new(model.User)

	u.ID = 1
	_, err := model.FindUserByUid(u, db.DB)
	assert.NoError(t, err)
}

func TestFindUserByUidError(t *testing.T) {
	u := new(model.User)

	u.ID = 999999
	_, err := model.FindUserByUid(u, db.DB)
	assert.Error(t, err, fmt.Sprint("該当のユーザーが見つかりません:record not found"))
}


/*
FindUserByMailPass
Normal
1. MailとPassに基づいてユーザー情報を取得

Error
1. 該当条件のユーザーが存在しない
*/
func TestFindUserByMailPassNormal(t *testing.T) {
	u := new(model.User)

	u.Password = "123"
	u.Mail = "sample1@gmail.com"
	_, err := model.FindUserByMailPass(u, db.DB)
	assert.NoError(t, err)
}

func TestFindUserByMailPassError(t *testing.T) {
	u := new(model.User)

	u.Password = "999999"
	u.Mail = "sample999@gmail.com"
	_, err := model.FindUserByMailPass(u, db.DB)
	assert.Error(t, err, fmt.Sprint("該当のユーザーが見つかりません:record not found"))
}
