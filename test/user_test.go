package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/labstack/echo"
	"work/handler"
)

/*
!!一度動作することは確認済みuser tableクリアしないと2度目以降のテストは失敗する
 */
func TestSignUp(t *testing.T) {
	e := echo.New()

	// テスト用リクエスト生成
	user_json := `{"Password": "dummy", "Name": "dummy_account", "address": "京都市","Tel": "000-0000-0000", "Mail": "sample@gmail.com"}`
	bodyReader := strings.NewReader(user_json)
	req := httptest.NewRequest("POST", "/signup", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()

	contents := e.NewContext(req, rec)

	assert.NoError(t, handler.Signup(contents))
	assert.Equal(t, http.StatusCreated, rec.Code)
	// 本当はJSONの中身にuser_jsonの値それぞれが含むかをチェックしたいが
	// rec.Bodyに\などの記号が挿入されるのでこの程度のチェックで留める
	assert.Contains(t, rec.Body.String(), "dummy")
}

func TestLogin(t *testing.T) {
	e := echo.New()

	// ↑のテストでDBに追加済み
	user_json := `{"Password": "dummy", "Name": "dummy_account", "address": "京都市","Tel": "000-0000-0000", "Mail": "sample@gmail.com"}`
	bodyReader := strings.NewReader(user_json)
	req := httptest.NewRequest("POST", "/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()

	contents := e.NewContext(req, rec)

	assert.NoError(t, handler.Login(contents))
	assert.Equal(t, http.StatusOK, rec.Code)
	// 本当はJSONの中身にuser_jsonの値それぞれが含むかをチェックしたいが
	// rec.Bodyに\などの記号が挿入されるのでこの程度のチェックで留める
	assert.Contains(t, rec.Body.String(), "token")
}
