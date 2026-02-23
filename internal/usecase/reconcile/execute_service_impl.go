package reconcile

import (
	"recon-service/internal/domain/reconcile"
	"recon-service/internal/shared/constants"
)

type matchKey struct {
	Date   string
	Type   constants.TransactionType
	Amount string
}

func (s *reconcileService) Execute(system []reconcile.SystemTransaction, banks []reconcile.BankTransaction) (reconcile.Summary, error) {
	index := make(map[matchKey][]*reconcile.SystemTransaction, len(system))

	for i := range system {
		trx := &system[i]

		key := matchKey{
			Date:   trx.TransactionTime.Format("2006-01-02"),
			Type:   trx.Type,
			Amount: trx.Amount,
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
			Date:   bankTrx.Date.Format("2006-01-02"),
			Type:   bankTrx.Type,
			Amount: bankTrx.Amount,
		}

		if list, ok := index[key]; ok && len(list) > 0 {

			summary.TotalMatched++

			// remove one element
			if len(list) == 1 {
				delete(index, key)
			} else {
				index[key] = list[1:]
			}

		} else {
			summary.UnmatchedBankByBank[bankTrx.BankName] =
				append(summary.UnmatchedBankByBank[bankTrx.BankName], bankTrx)
		}
	}

	// Remaining = unmatched system
	for _, remaining := range index {
		for _, trx := range remaining {
			summary.UnmatchedSystem = append(summary.UnmatchedSystem, *trx)
		}
	}

	summary.TotalUnmatched = len(summary.UnmatchedSystem)
	for _, v := range summary.UnmatchedBankByBank {
		summary.TotalUnmatched += len(v)
	}

	return summary, nil
}
