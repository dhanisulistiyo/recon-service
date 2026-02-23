package job

type Repository interface {
	Create(job *Job)
	GetByIdempotencyKey(key string) (*Job, bool)
	Update(job *Job)
	Get(id string) (*Job, bool)
}
