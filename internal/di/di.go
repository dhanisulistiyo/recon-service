package di

import (
	jobInfr "recon-service/internal/infrastructure/memory/job"
	httpHandler "recon-service/internal/interface/http"
	"recon-service/internal/usecase/job"
	"recon-service/internal/usecase/reconcile"
)

func Initialize() *httpHandler.Handler {
	jobRepo := jobInfr.NewInMemoryRepo()
	reconSvc := reconcile.NewService()
	jobSvc := job.NewService(jobRepo)
	workerSvc := job.NewWorker(jobRepo, reconSvc)
	handler := httpHandler.NewHandler(jobSvc, workerSvc)

	return handler
}
