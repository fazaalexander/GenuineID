package main

import (
	"github.com/fazaalexander/GenuineID/config"
	m "github.com/fazaalexander/GenuineID/middlewares"
	"github.com/fazaalexander/GenuineID/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	m.LogMiddleware(e)
	e.Logger.Fatal(e.Start(":8000"))
}
