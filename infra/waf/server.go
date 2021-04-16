package waf

import (
	"fmt"
	_ "github.com/SakagamiKazuto/golang_api/docs"
	"github.com/SakagamiKazuto/golang_api/infra/dbhandle"
	"github.com/SakagamiKazuto/golang_api/interface/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
	"os"
)

type Server struct {
	*echo.Echo
}

func createServer() (*Server, error) {
	return &Server{
		Echo: echo.New(),
	}, nil
}

func (s Server) setRouter() {
	e := s.Echo

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})

	s.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	{
		uc := controller.NewUserController(dbhandle.NewDBHandler())
		s.Echo.POST("/signup", s.SignUp(uc))
		s.Echo.POST("/login", s.Login(uc))
	}

	// JWTの認証を必要とするAPIは以下に記述
	{
		bc := controller.NewBosyuController(dbhandle.NewDBHandler())
		api := s.Echo.Group("/api")
		api.Use(middleware.JWTWithConfig(JwtConf))
		api.POST("/bosyu/create", s.CreateBosyu(bc))
		api.GET("/bosyu/get", s.GetBosyu(bc))
		api.PUT("/bosyu/update", s.UpdateBosyu(bc))
		api.DELETE("/bosyu/delete", s.DeleteBosyu(bc))
	}
}

func (s *Server) run() {
	s.Echo.Logger.Fatal(s.Echo.Start(":" + s.getPort()))
}

func (s *Server) getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	return port
}

func Run() {
	s, err := createServer()
	if err != nil {
		logrus.Fatal(fmt.Sprintf("サーバー起動時にエラーが発生しました\n%s", err.Error()))
		panic("処理を中断しました")
	}
	s.setRouter()
	s.run()
}
