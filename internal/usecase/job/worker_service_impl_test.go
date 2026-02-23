package job

import (
	"testing"

	jobEntity "recon-service/internal/domain/job"
	"recon-service/internal/domain/reconcile"
	"recon-service/internal/shared/constants"
	jobMocks "recon-service/mocks/domain/job"
	reconMocks "recon-service/mocks/usecase/reconcile"

	"github.com/stretchr/testify/mock"
)

func TestWorker_Process(t *testing.T) {
	repo := new(jobMocks.Repository)
	reconSvc := new(reconMocks.Service)
	worker := &Worker{
		repo:    repo,
		service: reconSvc,
		queue:   make(chan *JobPayload, 1),
	}

	t.Run("successful process", func(t *testing.T) {
		jobID := "job-1"
		payload := &JobPayload{
			JobID:  jobID,
			Start:  "2024-01-01",
			End:    "2024-01-31",
			System: []byte("trxID,amount,type,transactionTime\nTRX001,100.00,CREDIT,2024-01-10T10:00:00Z"),
			Banks: []BankFile{
				{
					Name: "BCA",
					Data: []byte("id,amount,date\nBCA001,100.00,2024-01-10"),
				},
			},
		}

		existingJob := &jobEntity.Job{ID: jobID, Status: constants.StatusProcessing}
		repo.On("Get", jobID).Return(existingJob, true).Once()

		reconSvc.On("Execute", mock.Anything, mock.Anything).Return(reconcile.Summary{
			TotalProcessed: 2,
			TotalMatched:   1,
		}, nil).Once()

		repo.On("Update", mock.MatchedBy(func(j *jobEntity.Job) bool {
			return j.ID == jobID && j.Status == constants.StatusDone && j.Result.TotalMatched == 1
		})).Return().Once()

		worker.process(payload)

		repo.AssertExpectations(t)
		reconSvc.AssertExpectations(t)
	})

	t.Run("job not found", func(t *testing.T) {
		jobID := "missing-job"
		payload := &JobPayload{JobID: jobID}

		repo.On("Get", jobID).Return(nil, false).Once()

		worker.process(payload)

		repo.AssertExpectations(t)
	})

	t.Run("invalid transactions - should fail job", func(t *testing.T) {
		jobID := "job-fail"
		payload := &JobPayload{
			JobID: jobID,
			Start: "invalid-date",
		}

		existingJob := &jobEntity.Job{ID: jobID, Status: constants.StatusProcessing}
		repo.On("Get", jobID).Return(existingJob, true).Once()

		repo.On("Update", mock.MatchedBy(func(j *jobEntity.Job) bool {
			return j.ID == jobID && j.Status == constants.StatusFailed
		})).Return().Once()

		worker.process(payload)

		repo.AssertExpectations(t)
	})
}
