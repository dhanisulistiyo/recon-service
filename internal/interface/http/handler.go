package http

import (
	"recon-service/internal/usecase/job"
)

type Handler struct {
	jobUsecase job.Service
	worker     *job.Worker
}

func NewHandler(jobUsecase job.Service, worker *job.Worker) *Handler {
	return &Handler{jobUsecase: jobUsecase, worker: worker}
}
