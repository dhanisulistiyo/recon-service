package http

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	jobUsecase "recon-service/internal/usecase/job"
)

func (h *Handler) Reconcile(c echo.Context) error {
	// read system file
	systemFile, err := c.FormFile("system_file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "system_file is required")
	}
	systemReader, err := systemFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open system_file")
	}
	defer systemReader.Close()
	systemBytes, err := io.ReadAll(systemReader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read system_file")
	}

	// read bank files
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid multipart form")
	}
	bankFiles := form.File["bank_files"]
	if len(bankFiles) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "at least one bank_file is required")
	}

	start, err := time.Parse("2006-01-02", c.FormValue("start_date"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start_date")
	}
	end, err := time.Parse("2006-01-02", c.FormValue("end_date"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid end_date")
	}

	var banks []jobUsecase.BankFile
	for _, fh := range bankFiles {
		f, err := fh.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to open bank_file: "+fh.Filename)
		}
		b, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read bank_file: "+fh.Filename)
		}
		// strip directory prefix and file extension (e.g. "bca.csv" → "bca")
		name := strings.TrimSuffix(filepath.Base(fh.Filename), filepath.Ext(fh.Filename))
		banks = append(banks, jobUsecase.BankFile{Name: name, Data: b})
	}

	// --- generate idempotency key ---
	hash := sha256.New()
	hash.Write(systemBytes)
	for _, b := range banks {
		hash.Write(b.Data)
	}
	hash.Write([]byte(start.Format("2006-01-02")))
	hash.Write([]byte(end.Format("2006-01-02")))
	idempotencyKey := hex.EncodeToString(hash.Sum(nil))

	// --- create new job and check idempotent ---
	jobEntity, existJob := h.jobUsecase.CreateWithKey(idempotencyKey)
	if existJob {
		return c.JSON(http.StatusAccepted, map[string]string{
			"job_id": jobEntity.ID,
		})
	}

	h.worker.Enqueue(&jobUsecase.JobPayload{
		JobID:  jobEntity.ID,
		System: systemBytes,
		Banks:  banks,
		Start:  start,
		End:    end,
	})

	return c.JSON(http.StatusAccepted, map[string]string{
		"job_id": jobEntity.ID,
	})
}

func (h *Handler) GetJob(c echo.Context) error {
	id := c.Param("id")

	job, err := h.jobUsecase.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}

	return c.JSON(http.StatusOK, job)
}
