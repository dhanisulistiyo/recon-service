package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/health", h.Health)
	e.POST("/reconcile", h.Reconcile)
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "reconciliation Service. App Version 1: "+time.Now().Format(time.RFC3339))
}
