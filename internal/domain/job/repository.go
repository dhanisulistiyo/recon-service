package job

type Repository interface {
	Create(job *Job)
	Update(job *Job)
	Get(id string) (*Job, bool)
}
