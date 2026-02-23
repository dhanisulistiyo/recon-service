package job

import (
	"time"

	"recon-service/internal/domain/reconcile"
	"recon-service/internal/shared/constants"
)

type Job struct {
	ID             string             `json:"id"`
	Status         constants.Status   `json:"status"`
	Result         *reconcile.Summary `json:"result"`
	Error          string             `json:"error"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	IdempotencyKey string             `json:"-"`
}
