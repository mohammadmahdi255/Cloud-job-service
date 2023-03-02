package handler

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Authentication(c echo.Context) error {
	return nil
}

func (h *Handler) Upload(c echo.Context) error {

	// todo: check token
	// todo: insert new row in DBaas
	// todo: upload file
	return nil
}

func (h *Handler) Execute(c echo.Context) error {
	return nil
}

func (h *Handler) JobStatus(c echo.Context) error {
	return nil
}
