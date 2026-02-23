package job

import (
	"errors"
	"time"

	"github.com/google/uuid"

	jobEntity "recon-service/internal/domain/job"
	"recon-service/internal/shared/constants"
)

func (s *service) Create() (*jobEntity.Job, error) {
	j := &jobEntity.Job{
		ID:        uuid.New().String(),
		Status:    constants.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.repo.Create(j)
	return j, nil
}

func (s *service) Get(id string) (*jobEntity.Job, error) {
	j, ok := s.repo.Get(id)
	if !ok {
		return nil, errors.New("job not found")
	}
	return j, nil
}
