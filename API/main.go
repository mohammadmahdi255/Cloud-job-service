package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/handler"
)

func main() {
	e := echo.New()
	g := e.Group("/api")

	h := handler.NewHandler()
	h.RegisterRoutes(g)

	e.Logger.Fatal(e.Start(":8080"))
}
