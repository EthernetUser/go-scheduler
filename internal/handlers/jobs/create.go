package handlers

import (
	jobsService "go-scheduler/internal/jobs"

	"github.com/gin-gonic/gin"
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

//	@BasePath	/api/v1

// @Summary	Add a new job
// @ID			create-job
// @Accept		json
// @Produce	json
// @Param		request	body		CreateJobRequest	true	"Some ID"
// @Success	200		{object}	CreateJobResponse	"ok"
// @Router		/jobs [post]
func (h *handlers) Create() gin.HandlerFunc {
	logger := h.logger.WithGroup("handler.create")
	return func(ctx *gin.Context) {
		var job CreateJobRequest
		if err := ctx.BindJSON(&job); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		parsedJob := jobsService.JobForCreation{
			Name:        job.Name,
			Description: job.Description,
			Cron:        job.Cron,
			Url:         job.Url,
		}

		err := h.jobService.CreateJob(ctx, parsedJob)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "success"})
	}
}
