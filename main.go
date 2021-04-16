package main

import (
	"github.com/SakagamiKazuto/golang_api/infra/waf"
	_ "github.com/SakagamiKazuto/golang_api/logger"
)

// @title matchihg_portfolio
// @version 1.0
// @description This is goecho api server.
// @host localhost:9999
// @BasePath /
func main() {
	//db.InitDB()
	waf.Run()
}


