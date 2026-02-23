package job

import (
	"testing"

	"recon-service/internal/domain/job"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepo(t *testing.T) {
	repo := NewInMemoryRepo()

	t.Run("Create and Get", func(t *testing.T) {
		j := &job.Job{ID: "job-1"}
		repo.Create(j)

		found, ok := repo.Get("job-1")
		assert.True(t, ok)
		assert.Equal(t, j, found)
	})

	t.Run("Update", func(t *testing.T) {
		j := &job.Job{ID: "job-1", Error: ""}
		repo.Create(j)

		j.Error = "some error"
		repo.Update(j)

		found, ok := repo.Get("job-1")
		assert.True(t, ok)
		assert.Equal(t, "some error", found.Error)
	})

	t.Run("Get missing", func(t *testing.T) {
		found, ok := repo.Get("missing")
		assert.False(t, ok)
		assert.Nil(t, found)
	})
}
