package reconcile

import "time"

type TransactionType string

const (
	Debit  TransactionType = "DEBIT"
	Credit TransactionType = "CREDIT"
)

type SystemTransaction struct {
	TrxID           string
	Amount          string
	Type            TransactionType
	TransactionTime time.Time
}

type BankTransaction struct {
	UniqueID string
	Amount   string
	Type     TransactionType
	Date     time.Time
	BankName string
}

type Summary struct {
	TotalProcessed      int
	TotalMatched        int
	TotalUnmatched      int
	TotalDiscrepancy    float64
	UnmatchedSystem     []SystemTransaction
	UnmatchedBankByBank map[string][]BankTransaction
}

/*
Repository abstraction
Mockable via mockery
*/
type Repository interface {
	SaveSummary(summary Summary) error
}
