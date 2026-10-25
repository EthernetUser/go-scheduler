package handlers

import (
	"encoding/json"
	jobsService "go-scheduler/internal/jobs"
	"net/http"
)

type CreateJobRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cron        string `json:"cron"`
	Url         string `json:"url"`
}

type CreateJobResponse struct {
	Message string `json:"message"`
}

// TODO: update swagger docs
//	@BasePath	/api/v1

// @Summary	Add a new job
// @ID			create-job
// @Accept		json
// @Produce	json
// @Param		request	body		CreateJobRequest	true	"Some ID"
// @Success	200		{object}	CreateJobResponse	"ok"
// @Router		/jobs [post]
func (h *handlers) Create() http.HandlerFunc {
	logger := h.logger.WithGroup("handler.create")
	return func(w http.ResponseWriter, r *http.Request) {
		var job CreateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			logger.Error("Failed to parse json body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedJob := jobsService.JobForCreation{
			Name:        job.Name,
			Description: job.Description,
			Cron:        job.Cron,
			Url:         job.Url,
		}

		err := h.jobService.CreateJob(r.Context(), parsedJob)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			//TODO: add checking err type everywhere to provide correct err status code
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
