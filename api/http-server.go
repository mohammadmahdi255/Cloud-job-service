package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	"github.com/mohammadmahdi255/Cloud-job-service/handler"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/services"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
	"os"
	"strings"
	"sync"
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
	s := storage.NewStorage("default", global.Endpoint)
	p := rabbitmq.NewProducer(url)

	err := s.CreateBucket(global.Bucket)
	if err != nil {
		fmt.Println(err)
	}

	objects, err := s.GetListObject(global.Bucket)
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

	mutex := &sync.RWMutex{}
	mutex.Lock()
	go func() {
		services.JobMaker(mutex, global.Bucket, global.Endpoint, url)
	}()

	mutex.Lock()
	go func() {
		services.ExecuteJob(mutex)
	}()

	apiBasePath := "/api/auth"
	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://dev-f6572cc1be7611edbd48135fc52bc1b6-ap-southeast-1.aws.supertokens.io:3573",
			APIKey:        "v7BAm-gbgWJ6xKtw9K=RGSTV9WcgD3",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "cloud-job-service",
			APIDomain:     "http://localhost:8080",
			APIBasePath:   &apiBasePath,
			WebsiteDomain: "http://localhost:8080",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired, // or models.ModeOptional
			}),
			session.Init(&sessmodels.TypeInput{}),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	// CORS middleware
	e.Use(func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			if c.Request().Method == "OPTIONS" {
				c.Response().Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
				c.Response().Header().Set("Access-Control-Allow-Methods", "*")
				_, err := c.Response().Write([]byte(""))
				if err != nil {
					return err
				}
				return nil
			} else {
				return hf(c)
			}
		}
	})

	// SuperTokens Middleware
	e.Use(func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				if err := hf(c); err != nil {
					c.Error(err)
				}
			})).ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})

	e.Logger.Fatal(e.Start("localhost:8080"))

}
