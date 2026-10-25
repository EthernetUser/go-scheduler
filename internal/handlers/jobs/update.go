package handlers

import (
	"encoding/json"
	jobsService "go-scheduler/internal/jobs"
	"net/http"
)

type UpdateJobRequest struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cron        string `json:"cron"`
	Url         string `json:"url"`
}

type UpdateJobResponse struct {
	Message string `json:"message"`
}

//	@BasePath	/api/v1

// @Summary	Update job
// @ID			update-job
// @Accept		json
// @Produce	json
// @Param		request	body		UpdateJobRequest	true	"Some ID"
// @Success	200		{object}	UpdateJobResponse	"ok"
// @Router		/jobs [put]
func (h *handlers) Update() http.HandlerFunc {
	logger := h.logger.WithGroup("handler.update")
	return func(w http.ResponseWriter, r *http.Request) {
		var job UpdateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			logger.Error("Failed to parse requet body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedJob := jobsService.Job{
			Id:          job.Id,
			Name:        job.Name,
			Description: job.Description,
			Cron:        job.Cron,
			Url:         job.Url,
		}

		err := h.jobService.UpdateJob(r.Context(), parsedJob)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
