package handlers

import (
	jobsService "go-scheduler/internal/jobs"

	"github.com/gin-gonic/gin"
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
func (h *handlers) Get() gin.HandlerFunc {
	logger := h.logger.WithGroup("handler.get")
	return func(ctx *gin.Context) {
		jobs, err := h.jobService.GetJobs(ctx)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"jobs": jobs})
	}
}
