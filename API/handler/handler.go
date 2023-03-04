package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/handler/request/models"
	ResponseModels "github.com/mohammadmahdi255/Cloud-job-service/handler/response/models"
	"io"
	"mime/multipart"
	"net/http"
)

type Handler struct {
	database *database.Database
}

func NewHandler(database *database.Database) *Handler {
	return &Handler{database}
}

func (h *Handler) Authentication(c echo.Context) error {
	return nil
}

func (h *Handler) Upload(c echo.Context) error {

	fileHeader, err := c.FormFile("programFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	uploadFile, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	for {
		p := make([]byte, 5)
		_, err := uploadFile.Read(p)
		if err == io.EOF {
			break;
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s", p)
	}

	fmt.Println()

	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			panic(err)
		}
	}(uploadFile)

	// todo: get json data
	upload := models.NewUpload(c.Request().MultipartForm.Value["json"])

	// todo: check token

	// todo: insert new row in DBaas
	err = h.database.AddUpload(upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseModels.NewMessage(err.Error()))
	}

	// todo: upload file
	return c.JSON(http.StatusOK, ResponseModels.NewMessage(upload))
}

func (h *Handler) Execute(c echo.Context) error {
	return nil
}

func (h *Handler) JobStatus(c echo.Context) error {
	return nil
}
