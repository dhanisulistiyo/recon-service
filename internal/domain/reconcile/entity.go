package reconcile

import (
	"time"

	"recon-service/internal/shared/constants"
)

type SystemTransaction struct {
	TrxID           string                    `json:"trx_id"`
	Amount          string                    `json:"amount"`
	Type            constants.TransactionType `json:"type"`
	TransactionTime time.Time                 `json:"transaction_time"`
}

type BankTransaction struct {
	UniqueID string                    `json:"unique_id"`
	Amount   string                    `json:"amount"`
	Type     constants.TransactionType `json:"type"`
	Date     time.Time                 `json:"date"`
	BankName string                    `json:"bank_name"`
}

type Summary struct {
	TotalProcessed      int                          `json:"total_processed"`
	TotalMatched        int                          `json:"total_matched"`
	TotalUnmatched      int                          `json:"total_unmatched"`
	TotalDiscrepancy    float64                      `json:"total_discrepancy"`
	UnmatchedSystem     []SystemTransaction          `json:"unmatched_system"`
	UnmatchedBankByBank map[string][]BankTransaction `json:"unmatched_bank_by_bank"`
}
