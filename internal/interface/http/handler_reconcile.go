package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Reconcile(c echo.Context) error {
	// In real implementation:
	// - parse multipart
	// - parse CSV
	// - map to domain entities
	// For brevity here we assume already mapped

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Reconciliation endpoint ready",
	})
}
