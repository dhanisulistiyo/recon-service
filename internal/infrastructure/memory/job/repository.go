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

func (r *inMemoryRepo) Create(j *job.Job) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// store a copy so callers can't mutate the stored value directly
	cp := *j
	r.data[j.ID] = &cp
}

func (r *inMemoryRepo) GetByIdempotencyKey(key string) (*job.Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, job := range r.data {
		if job.IdempotencyKey == key {
			return job, true
		}
	}
	return nil, false
}

func (r *inMemoryRepo) Update(j *job.Job) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// always store a fresh copy to prevent data races
	cp := *j
	r.data[j.ID] = &cp
}

// Get returns a copy of the stored job to prevent concurrent mutation.
func (r *inMemoryRepo) Get(id string) (*job.Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.data[id]
	if !ok {
		return nil, false
	}
	cp := *j
	return &cp, true
}
