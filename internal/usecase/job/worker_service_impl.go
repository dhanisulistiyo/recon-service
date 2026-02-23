package job

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"time"

	"recon-service/internal/domain/reconcile"
	"recon-service/internal/shared/constants"
)

func (w *Worker) Enqueue(payload *JobPayload) {
	w.queue <- payload
}

func (w *Worker) start() {
	for payload := range w.queue {
		w.process(payload)
	}
}

func (w *Worker) process(payload *JobPayload) {
	job, ok := w.repo.Get(payload.JobID)
	if !ok {
		return
	}

	systemTrx, bankTrx, err := w.prepareTransactions(payload)
	if err != nil {
		job.Status = constants.StatusFailed
		job.Error = err.Error()
		w.repo.Update(job)
		return
	}

	summary, err := w.service.Execute(systemTrx, bankTrx)
	if err != nil {
		job.Status = constants.StatusFailed
		job.Error = err.Error()
		w.repo.Update(job)
		return
	}

	job.Status = constants.StatusDone
	job.Result = &summary
	job.UpdatedAt = time.Now()

	w.repo.Update(job)
}

func (w *Worker) prepareTransactions(payload *JobPayload) ([]reconcile.SystemTransaction, []reconcile.BankTransaction, error) {
	start, err := time.Parse("2006-01-02", payload.Start)
	if err != nil {
		return nil, nil, err
	}

	end, err := time.Parse("2006-01-02", payload.End)
	if err != nil {
		return nil, nil, err
	}

	systemTrx, err := parseSystemCSVBytes(payload.System, start, end)
	if err != nil {
		return nil, nil, err
	}

	var allBankTrx []reconcile.BankTransaction

	for _, bankFile := range payload.Banks {
		trx, err := parseBankCSVBytes(bankFile.Data, bankFile.Name, start, end)
		if err != nil {
			return nil, nil, err
		}

		allBankTrx = append(allBankTrx, trx...)
	}

	return systemTrx, allBankTrx, nil
}

func parseSystemCSVBytes(data []byte, start, end time.Time) ([]reconcile.SystemTransaction, error) {

	reader := csv.NewReader(bytes.NewReader(data))

	// skip header
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var result []reconcile.SystemTransaction

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 4 {
			return nil, errors.New("invalid system csv format")
		}

		t, err := time.Parse(time.RFC3339, record[3])
		if err != nil || t.Before(start) || t.After(end) {
			continue
		}

		result = append(result, reconcile.SystemTransaction{
			TrxID:           record[0],
			Amount:          record[1],
			Type:            constants.TransactionType(record[2]),
			TransactionTime: t,
		})
	}

	return result, nil
}

func parseBankCSVBytes(data []byte, bankName string, start, end time.Time) ([]reconcile.BankTransaction, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	// skip header
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var result []reconcile.BankTransaction

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 3 {
			return nil, errors.New("invalid bank csv format")
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02", record[2])
		if err != nil {
			return nil, err
		}

		if date.Before(start) || date.After(end) {
			continue
		}

		var trxType constants.TransactionType
		if amount < 0 {
			trxType = constants.Debit
			amount = -amount
		} else {
			trxType = constants.Credit
		}

		result = append(result, reconcile.BankTransaction{
			UniqueID: record[0],
			Amount:   strconv.FormatFloat(amount, 'f', 2, 64),
			Type:     trxType,
			Date:     date,
			BankName: bankName,
		})
	}

	return result, nil
}
