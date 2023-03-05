package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/handler"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	_ "log"
	_ "net/http"
)

func main() {
	e := echo.New()
	g := e.Group("/api")

	endpoint := "https://s3.ir-thr-at1.arvanstorage.ir"

	d := database.NewDatabase()
	s := storage.NewStorage("default", endpoint)

	bucket := "cloud-job-service"

	err := s.CreateBucket(bucket)
	if err != nil {
		fmt.Println(err)
	}

	objects, err := s.GetListObject(bucket)
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range objects {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("ETag:", *item.ETag)
		fmt.Println("")
	}

	h := handler.NewHandler(d, s)
	h.RegisterRoutes(g)

	e.Logger.Fatal(e.Start(":8080"))

	//http.HandleFunc("/create", d.DatabaseCreate)
	//
	//log.Fatal(http.ListenAndServe(":8080", nil))

}
