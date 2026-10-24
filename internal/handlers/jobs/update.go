package handlers

import (
	jobsService "go-scheduler/internal/jobs"

	"github.com/gin-gonic/gin"
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
func (h *handlers) Update() gin.HandlerFunc {
	logger := h.logger.WithGroup("handler.update")
	return func(ctx *gin.Context) {
		var job UpdateJobRequest
		if err := ctx.BindJSON(&job); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		parsedJob := jobsService.Job{
			Id:          job.Id,
			Name:        job.Name,
			Description: job.Description,
			Cron:        job.Cron,
			Url:         job.Url,
		}

		err := h.jobService.UpdateJob(ctx, parsedJob)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "success"})
	}
}
