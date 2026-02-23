package job

import (
	"time"

	jobEntity "recon-service/internal/domain/job"
	"recon-service/internal/usecase/reconcile"
)

type Worker struct {
	repo    jobEntity.Repository
	service reconcile.Service
	queue   chan *JobPayload
}

type BankFile struct {
	Name string
	Data []byte
}

type JobPayload struct {
	JobID  string
	System []byte
	Banks  []BankFile
	Start  time.Time
	End    time.Time
}

func NewWorker(repo jobEntity.Repository, service reconcile.Service) *Worker {
	w := &Worker{
		repo:    repo,
		service: service,
		queue:   make(chan *JobPayload, 100),
	}
	go w.start()
	return w
}

type Service interface {
	Create() (*jobEntity.Job, error)
	CreateWithKey(idempotencyKey string) (res *jobEntity.Job, existJob bool)
	Get(id string) (*jobEntity.Job, error)
}

type service struct {
	repo jobEntity.Repository
}

func NewService(repo jobEntity.Repository) Service {
	return &service{repo: repo}
}
