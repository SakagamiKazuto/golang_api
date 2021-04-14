package test

import (
	"fmt"
	"net/http"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/handler"
)

/*
Bosyu
Create
*/
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

	mockReq := MockReq{bosyuJson, token, "/api/bosyu/create", "POST"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.CreateBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

/*
Bosyu
Get
*/
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

/*
Bosyu
Update
*/
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
	mockReq := MockReq{bosyuJson, token, "/api/bosyu/update", "PUT"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.UpdateBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

/*
Bosyu
Update
*/
func TestDeleteBosyuHandlerNormal(t *testing.T) {
	e := echo.New()

	token, err := createTokenFromSomeUser()
	if err != nil {
		t.Errorf("got error like: %+v", err)
	}

	mockReq := MockReq{`{"id": 3}`, token, "/api/bosyu/delete", "DELETE"}
	req, rec := mockReq.createReq()

	contents := e.NewContext(req, rec)
	exec := middleware.JWTWithConfig(handler.Config)(handler.DeleteBosyu)(contents)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

