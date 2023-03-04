package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/handler"

	_ "log"
	_ "net/http"
	//"github.com/mohammadmahdi255/Cloud-job-service/handler"
)

func main() {
	e := echo.New()
	g := e.Group("/api")

	d := database.NewDatabase()

	h := handler.NewHandler(d)
	h.RegisterRoutes(g)

	//http.HandleFunc("/create", d.DatabaseCreate)
	//
	//log.Fatal(http.ListenAndServe(":8080", nil))

	e.Logger.Fatal(e.Start(":8080"))
}
