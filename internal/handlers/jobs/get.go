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
func Get(jobs *jobsService.Jobs) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobs, err := jobs.GetJobs(ctx)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"jobs": jobs})
	}
}
