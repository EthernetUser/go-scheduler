package handlers

import (
	jobsService "go-scheduler/internal/jobs"

	"github.com/gin-gonic/gin"
)

type JobRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cron        string `json:"cron"`
	Url         string `json:"url"`
}

type JobResponse struct {
	Message string `json:"message"`
}

//	@BasePath	/api/v1

// @Summary	Add a new job
// @ID			create-job
// @Accept		json
// @Produce	json
// @Param		request	body		JobRequest	true	"Some ID"
// @Success	200		{object}	JobResponse	"ok"
// @Router		/jobs [post]
func CreateCreateHandler(jobs *jobsService.Jobs) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var job JobRequest
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

		err := jobs.CreateJob(ctx, parsedJob)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "success"})
	}
}
