package reconcile

import (
	"math"
	"recon-service/internal/domain/reconcile"
	"recon-service/internal/shared/constants"
	"strconv"
)

type matchKey struct {
	Date string
	Type constants.TransactionType
}

func (s *reconcileService) Execute(system []reconcile.SystemTransaction, banks []reconcile.BankTransaction) (reconcile.Summary, error) {
	// index system trx by date+type (ignore amount)
	index := make(map[matchKey][]*reconcile.SystemTransaction, len(system))
	for i := range system {
		trx := &system[i]
		key := matchKey{
			Date: trx.TransactionTime.Format("2006-01-02"),
			Type: trx.Type,
		}
		index[key] = append(index[key], trx)
	}

	summary := reconcile.Summary{
		UnmatchedBankByBank: make(map[string][]reconcile.BankTransaction),
	}

	summary.TotalProcessed = len(system) + len(banks)

	for i := range banks {
		bankTrx := banks[i]
		key := matchKey{
			Date: bankTrx.Date.Format("2006-01-02"),
			Type: bankTrx.Type,
		}

		if list, ok := index[key]; ok && len(list) > 0 {
			sysTrx := list[0]

			// remove one element
			if len(list) == 1 {
				delete(index, key)
			} else {
				index[key] = list[1:]
			}

			// convert amounts to float64
			sysAmt, err1 := strconv.ParseFloat(sysTrx.Amount, 64)
			bankAmt, err2 := strconv.ParseFloat(bankTrx.Amount, 64)
			if err1 == nil && err2 == nil {
				diff := math.Abs(sysAmt - bankAmt)
				if diff > 0 {
					summary.TotalDiscrepancy += diff
				}
			}

			summary.TotalMatched++
		} else {
			// unmatched bank
			summary.UnmatchedBankByBank[bankTrx.BankName] =
				append(summary.UnmatchedBankByBank[bankTrx.BankName], bankTrx)
		}
	}

	// remaining system trx = unmatched
	for _, remaining := range index {
		for _, trx := range remaining {
			summary.UnmatchedSystem = append(summary.UnmatchedSystem, *trx)
		}
	}

	// total unmatched
	summary.TotalUnmatched = len(summary.UnmatchedSystem)
	for _, v := range summary.UnmatchedBankByBank {
		summary.TotalUnmatched += len(v)
	}

	return summary, nil
}
