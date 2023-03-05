package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(g *echo.Group) {

	g.POST("/authentication", h.Authentication)
	g.POST("/upload", h.Upload)
	g.POST("/execute", h.Execute)

	g.GET("/job/status", h.JobStatus)

}
