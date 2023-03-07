package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/handler"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/services"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	_ "log"
	_ "net/http"
	"os"
)

const (
	ENDPOINT = "https://s3.ir-thr-at1.arvanstorage.ir"
	BUCKET   = "cloud-job-service"
)

func main() {
	e := echo.New()
	g := e.Group("/api")

	url := os.Getenv("CLOUDAMQP_URL")

	if url == "" {
		fmt.Println("localhost rabbitMQ")
		url = "amqp://guest:guest@localhost:5672/"
	}

	d := database.NewDatabase()
	s := storage.NewStorage("default", ENDPOINT)
	p := rabbitmq.NewProducer(url)

	err := s.CreateBucket(BUCKET)
	if err != nil {
		fmt.Println(err)
	}

	objects, err := s.GetListObject(BUCKET)
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

	h := handler.NewHandler(d, s, p)
	h.RegisterRoutes(g)

	go func() {
		services.JobMaker(BUCKET, ENDPOINT, url)
	}()

	go func() {
		services.ExecuteJob()
	}()

	e.Logger.Fatal(e.Start(":8080"))

}
