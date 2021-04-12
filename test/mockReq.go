package test

import (
	"fmt"
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

	if mr.method == "POST" || mr.method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	}

	if mr.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", mr.token))
	}
	rec := httptest.NewRecorder()
	return req, rec
}

