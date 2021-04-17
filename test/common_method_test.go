package test

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/domain"
	"github.com/SakagamiKazuto/golang_api/infra/waf"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
)

type MockReq struct {
	jsonStr string
	token   string
	target  string
	method  string
}

func (mr MockReq) createReq() (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(mr.method, mr.target, strings.NewReader(mr.jsonStr))
	req.Header.Add("Content-Type", "application/json")

	if mr.method == "POST" || mr.method == "PUT" {
		req.Header.Add("Accept", "application/json")
	}

	if mr.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", mr.token))
	}
	rec := httptest.NewRecorder()
	return req, rec
}

func createTokenFromSomeUser() (string, error) {
	user, err := ur.FindUserByUid(&domain.User{Model: gorm.Model{ID: 1}})
	if err != nil {
		panic(fmt.Sprintf(`テスト中にエラーが発生:%s`, err.Error()))
	}
	token, err := waf.CreateToken(user.ID, user.Mail)
	return token, err
}
