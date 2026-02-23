package job

import (
	mocks "recon-service/mocks/domain/job"
	"testing"

	jobEntity "recon-service/internal/domain/job"
	mockss "recon-service/mocks/domain/job"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJobService_Create(t *testing.T) {
	repo := new(mockss.Repository)
	s := NewService(repo)

	repo.On("Create", mock.Anything).Return()

	j, err := s.Create()

	assert.NoError(t, err)
	assert.NotEmpty(t, j.ID)
	assert.Equal(t, "PROCESSING", string(j.Status))
	repo.AssertExpectations(t)
}

func TestJobService_Get(t *testing.T) {
	repo := new(mocks.Repository)
	s := NewService(repo)

	t.Run("success", func(t *testing.T) {
		expectedJob := &jobEntity.Job{ID: "123"}
		repo.On("Get", "123").Return(expectedJob, true).Once()

		j, err := s.Get("123")

		assert.NoError(t, err)
		assert.Equal(t, expectedJob, j)
	})

	t.Run("not found", func(t *testing.T) {
		repo.On("Get", "456").Return(nil, false).Once()

		j, err := s.Get("456")

		assert.Error(t, err)
		assert.Nil(t, j)
		assert.Equal(t, "job not found", err.Error())
	})

	repo.AssertExpectations(t)
}
