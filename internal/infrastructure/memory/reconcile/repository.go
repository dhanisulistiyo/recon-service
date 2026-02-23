package reconcile

import reconcileRepo "recon-service/internal/domain/reconcile"

type container struct{}

type inMemoryRepo struct{}

func (r *inMemoryRepo) SaveSummary(summary reconcileRepo.Summary) error {
	return nil
}
