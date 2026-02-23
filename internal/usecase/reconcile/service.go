package reconcile

import "recon-service/internal/domain/reconcile"

type reconcileService struct {
}

func NewService() Service {
	return &reconcileService{}
}

type Service interface {
	Execute(system []reconcile.SystemTransaction, banks []reconcile.BankTransaction) (reconcile.Summary, error)
}
