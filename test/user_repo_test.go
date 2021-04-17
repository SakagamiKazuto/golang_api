package test

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/infra/dbhandle"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ur = database.UserRepository{dbhandle.DBHandle{NewTestDB()}}

/*
CreateUser:
Normal
1. データを作成する

Error
1. Mailの重複キー制約に違反
2. 内部エラー
*/
func TestCreateUserNormal(t *testing.T) {
	u := new(domain.User)
	u.Mail = "sample9@gmail.com"
	u.Password = "123123"
	_, err := ur.CreateUser(u)
	assert.NoError(t, err)
}

func TestCreateUserError(t *testing.T) {
	u := new(domain.User)
	u.Mail = "sample1@gmail.com"
	_, err := ur.CreateUser(u)
	assert.Error(t, err, fmt.Sprintf(`メールアドレス%sのデータ挿入に失敗しました:pq: duplicate key value violates unique constraint "users_mail_key"`, u.Mail))
}

/*
FindUserByUid
Normal
1. idに基づいてユーザー情報を取得

Error
1. 該当idのユーザーが存在しない
*/
func TestFindUserByUidNormal(t *testing.T) {
	u := new(domain.User)

	u.ID = 1
	_, err := ur.FindUserByUid(u)
	assert.NoError(t, err)
}

func TestFindUserByUidError(t *testing.T) {
	u := new(domain.User)

	u.ID = 999999
	_, err := ur.FindUserByUid(u)
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
	u := new(domain.User)

	u.Password = "123"
	u.Mail = "sample1@gmail.com"
	_, err := ur.FindUserByMailPass(u)
	assert.NoError(t, err)
}

func TestFindUserByMailPassError(t *testing.T) {
	u := new(domain.User)

	u.Password = "999999"
	u.Mail = "sample999@gmail.com"
	_, err := ur.FindUserByMailPass(u)
	assert.Error(t, err, fmt.Sprint("該当のユーザーが見つかりません:record not found"))
}
