package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/SakagamiKazuto/golang_api/handler"
)

/*
Handler
createUser
*/
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

/*
Handler
Login
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
