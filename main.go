package main

import "os"

// @title matchihg_portfolio
// @version 1.0
// @description This is goecho api server.
// @host localhost:9999
// @BasePath /
func main() {
	router := newRouter()
	router.Logger.Fatal(router.Start(":" + getPort()))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	return port
}
