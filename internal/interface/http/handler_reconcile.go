package http

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"

	jobUsecase "recon-service/internal/usecase/job"
)

func (h *Handler) Reconcile(c echo.Context) error {
	jobEntity, err := h.jobUsecase.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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

	h.worker.Enqueue(&jobUsecase.JobPayload{
		JobID:  jobEntity.ID,
		System: systemBytes,
		Banks:  banks,
		Start:  c.FormValue("start_date"),
		End:    c.FormValue("end_date"),
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
