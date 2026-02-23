package di

import (
	httpHandler "recon-service/internal/interface/http"
	"recon-service/internal/usecase/reconcile"
)

func Initialize() *httpHandler.Handler {
	// repo := &inMemoryRepo{}
	service := reconcile.NewService()
	handler := httpHandler.NewHandler(service)

	return handler
}
