package di

import (
	"reconciliation-service/internal/domain/reconciliation"
	httpHandler "reconciliation-service/internal/interface/http"
	"reconciliation-service/internal/usecase/reconcile"
)

type container struct{}

type inMemoryRepo struct{}

func (r *inMemoryRepo) SaveSummary(summary reconciliation.Summary) error {
	return nil
}

func Initialize() *httpHandler.Handler {
	repo := &inMemoryRepo{}
	service := reconcile.NewService(repo)
	handler := httpHandler.NewHandler(service)

	return handler
}
