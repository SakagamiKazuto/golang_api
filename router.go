package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
	_ "work/docs"

	"work/handler"
)

func newRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	// JWTの認証を必要とするAPIは以下に記述
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(handler.Config))
	api.GET("/bosyu/get", handler.GetBosyu)
	api.POST("/bosyu/create", handler.CreateBosyu)
	api.PUT("/bosyu/update", handler.UpdateBosyu)
	api.DELETE("/bosyu/delete", handler.DeleteBosyu)

	return e
}
