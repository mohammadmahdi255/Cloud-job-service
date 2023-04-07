package handler

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	ResponseModels "github.com/mohammadmahdi255/Cloud-job-service/handler/response/models"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	"github.com/stretchr/objx"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
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

	// todo: check session exist or not
	err := h.checkSessionValid(c)
	if err != nil {
		return err
	}

	sessionContainer := c.Get("session").(sessmodels.SessionContainer)
	userInfo, err := thirdpartyemailpassword.GetUserById(sessionContainer.GetUserID())
	if err != nil {
		// TODO: Handle error
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err))
	}

	// todo: get json data
	upload, err := models.NewUpload(c.FormValue("json"))
	upload.Email = userInfo.Email
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
	fileHeader, err := c.FormFile("programFile")
	if err != nil {
		h.database.Delete(&upload)
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	uploadFile, err := fileHeader.Open()
	if err != nil {
		h.database.Delete(&upload)
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	filename := fmt.Sprintf("%d.%s", upload.Id, upload.ProgramLanguage)
	err = h.store.Upload(global.Bucket, filename, uploadFile)
	if err != nil {
		h.database.Delete(&upload)
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	err = uploadFile.Close()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	// todo: put tag to track file names
	err = h.store.PutObjectTag(global.Bucket, filename, "Name", fileHeader.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	fmt.Println(ResponseModels.NewMessage(upload))
	return c.JSON(http.StatusOK, ResponseModels.NewMessage(upload))
}

func (h *Handler) Execute(c echo.Context) error {

	// todo: check session exist or not
	err := h.checkSessionValid(c)
	if err != nil {
		return err
	}

	sessionContainer := c.Get("session").(sessmodels.SessionContainer)
	userInfo, err := thirdpartyemailpassword.GetUserById(sessionContainer.GetUserID())
	if err != nil {
		// TODO: Handle error
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err))
	}

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

	if upload.Email != userInfo.Email {
		return c.JSON(http.StatusBadRequest, "Your no owner of this id")
	}

	// todo: check can be Executed
	if !upload.IsEnable {
		return c.JSON(http.StatusOK, ResponseModels.NewMessage("can not be Execute because enable is 0"))
	}

	// todo: send id with rabbitMQ to mail-service
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
	// todo: check session exist or not
	err := h.checkSessionValid(c)
	if err != nil {
		return err
	}

	sessionContainer := c.Get("session").(sessmodels.SessionContainer)
	userInfo, err := thirdpartyemailpassword.GetUserById(sessionContainer.GetUserID())
	if err != nil {
		// TODO: Handle error
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err))
	}

	result, err := h.database.GetAllUserResult(userInfo.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) SessionInfo(c echo.Context) error {

	// todo: check session exist or not
	err := h.checkSessionValid(c)
	if err != nil {
		return err
	}

	sessionContainer := c.Get("session").(sessmodels.SessionContainer)
	sessionData, _ := sessionContainer.GetSessionData()

	userInfo, err := thirdpartyemailpassword.GetUserById(sessionContainer.GetUserID())
	if err != nil {
		// TODO: Handle error
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err))
	}

	data := map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
		"userInfo":           userInfo,
	}
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) RevokeToken(c echo.Context) error {

	// todo: check session exist or not
	err := h.checkSessionValid(c)
	if err != nil {
		return err
	}

	// retrieve the session object as shown below
	sessionContainer := c.Get("session").(sessmodels.SessionContainer)

	// This will delete the session from the db and from the frontend (cookies)
	err = sessionContainer.RevokeSession()
	if err != nil {
		err2 := supertokens.ErrorHandler(err, c.Request(), c.Response())
		if err2 != nil {
			// TODO: Send 500 status code to client
			return err2
		}
		return err
	}

	data := map[string]interface{}{
		"Status": "ok",
	}

	return c.JSON(http.StatusOK, data)
}

func (h *Handler) checkSessionValid(c echo.Context) error {
	// retrieve the session object as shown below
	sessionContainer := c.Get("session").(sessmodels.SessionContainer)

	if sessionContainer == nil {
		return errors.New("no session found")
	}

	_, err := sessionContainer.GetSessionData()
	if err != nil {
		err2 := supertokens.ErrorHandler(err, c.Request(), c.Response())
		if err2 != nil {
			// TODO: Send 500 status code to client
			return err2
		}
		return err
	}

	return nil
}

func (h *Handler) ListUploads(c echo.Context) error {
	// todo: check session exist or not
	err := h.checkSessionValid(c)
	if err != nil {
		return err
	}

	sessionContainer := c.Get("session").(sessmodels.SessionContainer)
	userInfo, err := thirdpartyemailpassword.GetUserById(sessionContainer.GetUserID())
	if err != nil {
		// TODO: Handle error
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err))
	}

	uploads, err := h.database.GetAllUserUpload(userInfo.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseModels.NewMessage(err.Error()))
	}

	list, err := h.store.GetListObject(global.Bucket)

	arr := make([]map[string]interface{}, 0, 10)

	j := 0
	for i := 0; i < len(list); i++ {
		key := fmt.Sprintf("%d.%s", uploads[j]["id"], uploads[j]["program_language"])
		if *list[i].Key == key {
			dic := structs.Map(list[i])
			dic["upload"] = uploads[j]
			arr = append(arr, dic)
			j++
		}
	}

	return c.JSON(http.StatusOK, arr)
}
