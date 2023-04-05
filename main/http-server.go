package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	"github.com/mohammadmahdi255/Cloud-job-service/handler"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
	"strings"
)

func main() {
	e := echo.New()
	g := e.Group("/api")

	d := database.NewDatabase()
	s := storage.NewStorage("default", global.Endpoint)
	p := rabbitmq.NewProducer(global.CloudamqpUrl)

	err := s.CreateBucket(global.Bucket)
	if err != nil {
		fmt.Println(err)
	}

	h := handler.NewHandler(d, s, p)
	h.RegisterRoutes(g)

	apiBasePath := "/api/auth"
	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: global.URI,
			APIKey:        global.APIKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "cloud-job-service",
			APIDomain:     "http://localhost:8080",
			APIBasePath:   &apiBasePath,
			WebsiteDomain: "http://localhost:8080",
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{}),
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
