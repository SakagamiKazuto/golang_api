package test

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/infra/dbhandle"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/domain"
)

var br = database.BosyuRepository{dbhandle.DBHandle{NewTestDB()}}

/*
CreateBosyu:
Normal
1. 募集データが挿入できる
Error
1. internalエラーをリターンする
*/
func TestCreateBosyuNormal(t *testing.T) {
	b := new(domain.Bosyu)

	b.Title = "sample_title"
	b.About = "sample_about"
	b.Prefecture = "愛媛県"
	b.City = "松山市"
	b.Level = "player"
	b.UserID = 123123

	_, err := br.CreateBosyu(b)
	assert.NoError(t, err)
}

func TestCreateBosyuError(t *testing.T) {
	b := new(domain.Bosyu)

	// すでにID=1のデータが存在するため一意制約に違反する
	b.ID = 1
	_, err := br.CreateBosyu(b)
	assert.Error(t, err)
}

/*
FindBosyuByUserID:
Normal
1. ユーザーIDでデータが取れる

Error
1. 該当ユーザーIDが存在しない
2. 内部エラー
*/
func TestFindBosyuByUidNormal(t *testing.T) {
	b := new(domain.Bosyu)

	b.UserID = 1
	_, err := br.FindBosyuByUid(b.UserID)
	assert.NoError(t, err)
}

func TestFindBosyuByUidError(t *testing.T) {
	b := new(domain.Bosyu)

	b.UserID = 999999
	_, err := br.FindBosyuByUid(b.UserID)
	assert.Error(t, err, fmt.Sprintf("該当のユーザーID%dの募集は見つかりません:record not found", b.UserID))
}

/*
UpdateBosyu:
Normal
1. BosyuのIDがdatabaseに存在する

Error
1. BosyuのIDがdatabaseに存在しない
*/
func TestUpdateBosyuNormal(t *testing.T) {
	b := new(domain.Bosyu)

	b.Title = "sample1_title_updated"
	b.About = "sample1_about_updated"
	b.Prefecture = "sample1_pref_updated"
	b.City = "sample1_city_updated"
	b.Level = "sample1_level_updated"
	b.UserID = 1
	b.ID = 1
	bosyu, err := br.UpdateBosyu(b)
	if assert.NoError(t, err) {
		assert.Empty(t, nil, bosyu.DeletedAt)
		assert.Equal(t, "sample1_title_updated", bosyu.Title)
	}
}

func TestUpdateBosyuError(t *testing.T) {
	b := new(domain.Bosyu)

	b.Title = "sample1_title_updated"
	b.About = "sample1_about_updated"
	b.Prefecture = "sample1_pref_updated"
	b.City = "sample1_city_updated"
	b.Level = "sample1_level_updated"
	b.UserID = 1
	b.ID = 99999

	_, err := br.UpdateBosyu(b)
	assert.Error(t, err)
}

/*
DeleteBosyuTests
Model:
Normal
1. BosyuのIDがdatabaseに存在する

Error
1. BosyuのIDがdatabaseに存在しない
*/
func TestDeleteBosyuNormal(t *testing.T) {
	b := new(domain.Bosyu)
	b.ID = 1
	err := br.DeleteBosyu(b.ID)
	assert.NoError(t, err)
}

func TestDeleteBosyuError(t *testing.T) {
	b := new(domain.Bosyu)
	b.ID = 0
	err := br.DeleteBosyu(b.ID)
	assert.Error(t, err)
}
