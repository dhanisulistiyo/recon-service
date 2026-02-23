package http

import (
	"recon-service/internal/usecase/reconcile"
)

type Handler struct {
	service reconcile.Service
}

func NewHandler(service reconcile.Service) *Handler {
	return &Handler{service: service}
}
