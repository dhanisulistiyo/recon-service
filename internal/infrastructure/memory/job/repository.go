package job

import (
	"sync"

	"recon-service/internal/domain/job"
)

type inMemoryRepo struct {
	mu   sync.RWMutex
	data map[string]*job.Job
}

func NewInMemoryRepo() job.Repository {
	return &inMemoryRepo{
		data: make(map[string]*job.Job),
	}
}

func (r *inMemoryRepo) Create(job *job.Job) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[job.ID] = job
}

func (r *inMemoryRepo) Update(job *job.Job) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[job.ID] = job
}

func (r *inMemoryRepo) Get(id string) (*job.Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.data[id]
	return job, ok
}
