package handlers

import (
	"net/http"
	"strconv"
)

const JobIdPathParam string = "jobId"

//	@BasePath	/api/v1

// @Summary	delete a new job
// @ID			delete-job
// @Accept		json
// @Produce	json
// @Param		jobId	path		string	true	"Some ID"
// @Success	200		{string}	string	"ok"
// @Router		/jobs/{jobId} [delete]
func (h *handlers) Delete() http.HandlerFunc {
	logger := h.logger.WithGroup("handler.delete")
	return func(w http.ResponseWriter, r *http.Request) {
		jobId := r.PathValue(JobIdPathParam)
		if jobId == "" {
			logger.Error("Empty jobId")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedJobId, err := strconv.Atoi(jobId)
		if err != nil {
			logger.Error("jobId is not a number")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.jobService.DeleteJob(r.Context(), parsedJobId)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
