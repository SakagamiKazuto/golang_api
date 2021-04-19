package main

import (
	"github.com/SakagamiKazuto/golang_api/infra/waf"
	_ "github.com/SakagamiKazuto/golang_api/infra/waf/logger"
)

// @title golang_api
// @version 1.0
// @description This is goecho api server.
// @host localhost:9999
// @BasePath /
func main() {
	waf.Run()
}


