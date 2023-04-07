package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"net/http"
)

func (h *Handler) RegisterRoutes(g *echo.Group) {

	g.POST("/upload", h.Upload, h.verifySession(nil))
	g.POST("/execute", h.Execute, h.verifySession(nil))
	g.POST("/signout", h.RevokeToken, h.verifySession(nil))

	g.GET("/job/status", h.JobStatus, h.verifySession(nil))
	g.GET("/upload/info", h.ListUploads, h.verifySession(nil))
	g.GET("/session/info", h.SessionInfo, h.verifySession(nil))

}

func (h *Handler) verifySession(options *sessmodels.VerifySessionOptions) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
				c.Set("session", session.GetSessionFromRequestContext(r.Context()))
				err := hf(c)
				if err != nil {
					return
				}
			})(c.Response(), c.Request())
			return nil
		}
	}
}
