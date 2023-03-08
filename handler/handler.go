package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	ResponseModels "github.com/mohammadmahdi255/Cloud-job-service/handler/response/models"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	"github.com/stretchr/objx"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
)

type Handler struct {
	database *database.Database
	store    *storage.Storage
	producer *rabbitmq.Producer
}

func NewHandler(database *database.Database, store *storage.Storage, producer *rabbitmq.Producer) *Handler {
	return &Handler{database, store, producer}
}

func (h *Handler) Upload(c echo.Context) error {
	// todo: get json data
	upload, err := models.NewUpload(c.FormValue("json"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseModels.NewMessage(err.Error()))
	}

	// todo: insert new row in DBaas
	err = h.database.AddUpload(upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseModels.NewMessage(err.Error()))
	}

	fmt.Println(upload.Id)

	// todo: upload file in s3
	bucket := "cloud-job-service"
	fileHeader, err := c.FormFile("programFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	uploadFile, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	filename := fmt.Sprintf("%d.%s", upload.Id, upload.ProgramLanguage)
	err = h.store.Upload(bucket, filename, uploadFile)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	err = uploadFile.Close()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	// todo: put tag to track file names
	err = h.store.PutObjectTag(bucket, filename, fileHeader.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseModels.NewMessage(upload))
}

func (h *Handler) Execute(c echo.Context) error {

	// todo: get json data
	dic, err := objx.FromJSON(c.FormValue("json"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseModels.NewMessage(err.Error()))
	}

	// todo: get upload row
	upload, err := h.database.GetUpload(dic.Get("id").Int())
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseModels.NewMessage(err.Error()))
	}

	// todo: check can be Executed
	if !upload.IsEnable {
		return c.JSON(http.StatusOK, ResponseModels.NewMessage("can not be Execute because enable is 0"))
	}

	// todo: send id with rabbitMQ to services
	dic = objx.New(map[string]interface{}{
		"id": upload.Id,
	})
	err = h.producer.Send(dic.MustJSON())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseModels.NewMessage(fmt.Sprintf("new job created for upload with id: %d", upload.Id)))
}

func (h *Handler) JobStatus(c echo.Context) error {
	result, err := h.database.GetAllUserResult("adel110@aut.ac.ir")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) sessionInfo(c echo.Context) error {
	sessionContainer := c.Get("session").(sessmodels.SessionContainer)

	if sessionContainer == nil {
		return errors.New("no session found")
	}
	sessionData, err := sessionContainer.GetSessionData()
	if err != nil {
		err = supertokens.ErrorHandler(err, c.Request(), c.Response())
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	}

	data := map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	}
	return c.JSON(http.StatusOK, data)
}
