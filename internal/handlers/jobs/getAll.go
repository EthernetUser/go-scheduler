package handlers

import (
	"encoding/json"
	jobsService "go-scheduler/internal/jobs"
	"net/http"
)

type GetJobsResponse struct {
	Jobs []jobsService.Job `json:"jobs"`
}

//	@BasePath	/api/v1

// @Summary	get jobs
// @ID			get-jobs
// @Accept		json
// @Produce	json
// @Success	200	{object}	GetJobsResponse	"ok"
// @Router		/jobs [get]
func (h *handlers) GetAll() http.HandlerFunc {
	logger := h.logger.WithGroup("handler.get")
	return func(w http.ResponseWriter, r *http.Request) {
		jobs, err := h.jobService.GetJobs(r.Context())
		if err != nil {
			logger.Error("Error occurred", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(jobs)
		if err != nil {
			logger.Error("Failed to write jobs to response", "error", err)
		}
	}
}
