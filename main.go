package main

import (
	"github.com/fazaalexander/GenuineID/config"
	m "github.com/fazaalexander/GenuineID/middlewares"
	"github.com/fazaalexander/GenuineID/routes"
	"github.com/fazaalexander/GenuineID/seeders"
)

func main() {
	config.InitDB()
	seeders.ProductTypeSeed()
	e := routes.New()
	m.LogMiddleware(e)
	e.Logger.Fatal(e.Start(":8000"))
}
