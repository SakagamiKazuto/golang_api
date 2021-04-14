package test

import (
	"fmt"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/db"
	"github.com/SakagamiKazuto/golang_api/model"
)

/*
CreateBosyu:
Normal
1. 募集データが挿入できる
Error
1. internalエラーをリターンする
*/
func TestCreateBosyuNormal(t *testing.T) {
	b := new(model.Bosyu)

	b.Title = "sample_title"
	b.About = "sample_about"
	b.Prefecture = "愛媛県"
	b.City = "松山市"
	b.Level = "player"
	b.UserID = 123123

	_, err := model.CreateBosyu(b, db.DB)
	assert.NoError(t, err)
}

func TestCreateBosyuError(t *testing.T) {
	b := new(model.Bosyu)

	// すでにID=1のデータが存在するため一意制約に違反する
	b.ID = 1
	_, err := model.CreateBosyu(b, db.DB)
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
	b := new(model.Bosyu)

	b.UserID = 1
	_, err := model.FindBosyuByUid(b.UserID, db.DB)
	assert.NoError(t, err)
}

func TestFindBosyuByUidError(t *testing.T) {
	b := new(model.Bosyu)

	b.UserID = 999999
	_, err := model.FindBosyuByUid(b.UserID, db.DB)
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

func TestUpdateBosyuError(t *testing.T) {
	b := new(model.Bosyu)

	b.Title = "sample1_title_updated"
	b.About = "sample1_about_updated"
	b.Prefecture = "sample1_pref_updated"
	b.City = "sample1_city_updated"
	b.Level = "sample1_level_updated"
	b.UserID = 1
	b.ID = 99999

	_, err := model.UpdateBosyu(b, db.DB)
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
	b := new(model.Bosyu)
	b.ID = 1
	err := model.DeleteBosyu(b.ID, db.DB)
	assert.NoError(t, err)
}

func TestDeleteBosyuError(t *testing.T) {
	b := new(model.Bosyu)
	b.ID = 0
	err := model.DeleteBosyu(b.ID, db.DB)
	assert.Error(t, err)
}

