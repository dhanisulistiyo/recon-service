package reconcile

import (
	"testing"
	"time"

	"recon-service/internal/domain/reconcile"
	"recon-service/internal/shared/constants"

	"github.com/stretchr/testify/assert"
)

func TestReconcileService_Execute(t *testing.T) {
	s := NewService()

	t.Run("should match transactions with same date, type and amount", func(t *testing.T) {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

		system := []reconcile.SystemTransaction{
			{
				TrxID:           "TRX001",
				Amount:          "100000.00",
				Type:            constants.Credit,
				TransactionTime: date.Add(time.Hour * 10),
			},
		}

		banks := []reconcile.BankTransaction{
			{
				UniqueID: "BCA001",
				Amount:   "100000.00",
				Type:     constants.Credit,
				Date:     date,
				BankName: "BCA",
			},
		}

		summary, err := s.Execute(system, banks)

		assert.NoError(t, err)
		assert.Equal(t, 2, summary.TotalProcessed)
		assert.Equal(t, 1, summary.TotalMatched)
		assert.Equal(t, 0, summary.TotalUnmatched)
		assert.Len(t, summary.UnmatchedSystem, 0)
		assert.Len(t, summary.UnmatchedBankByBank, 0)
	})

	t.Run("should handle unmatched system transaction", func(t *testing.T) {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

		system := []reconcile.SystemTransaction{
			{
				TrxID:           "TRX001",
				Amount:          "100000.00",
				Type:            constants.Credit,
				TransactionTime: date,
			},
		}

		banks := []reconcile.BankTransaction{}

		summary, err := s.Execute(system, banks)

		assert.NoError(t, err)
		assert.Equal(t, 1, summary.TotalProcessed)
		assert.Equal(t, 0, summary.TotalMatched)
		assert.Equal(t, 1, summary.TotalUnmatched)
		assert.Len(t, summary.UnmatchedSystem, 1)
		assert.Equal(t, "TRX001", summary.UnmatchedSystem[0].TrxID)
	})

	t.Run("should handle unmatched bank transaction", func(t *testing.T) {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

		system := []reconcile.SystemTransaction{}

		banks := []reconcile.BankTransaction{
			{
				UniqueID: "BCA001",
				Amount:   "100000.00",
				Type:     constants.Credit,
				Date:     date,
				BankName: "BCA",
			},
		}

		summary, err := s.Execute(system, banks)

		assert.NoError(t, err)
		assert.Equal(t, 1, summary.TotalProcessed)
		assert.Equal(t, 0, summary.TotalMatched)
		assert.Equal(t, 1, summary.TotalUnmatched)
		assert.Len(t, summary.UnmatchedBankByBank["BCA"], 1)
		assert.Equal(t, "BCA001", summary.UnmatchedBankByBank["BCA"][0].UniqueID)
	})
}
