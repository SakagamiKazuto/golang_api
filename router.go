package main

import (
	//"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
	_ "work/docs"
	//"github.com/SakagamiKazuto/matching_portfolio/handler"
	"work/handler"
)

func newRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/signup", handler.Signup)       // POST /signup
	e.POST("/login", handler.Login)         // POST /login

	api := e.Group("/api")

	// apiテストにHandler定義してmockDB生成する方針でやった結果この記述が必要になった
	api.Use(middleware.JWTWithConfig(handler.Config))   // /api 下はJWTの認証が必要
	api.GET("/bosyu/get", handler.GetBosyu)
	api.POST("/bosyu/create", handler.CreateBosyu)
	//e.PUT("/bosyu/update", handler.UpdateBosyu)
	//e.DELETE("/bosyu/delete", handler.DeleteBosyu)

	return e
}
