package get

import (
	jobsService "go-scheduler/internal/domain/jobs"

	"github.com/gin-gonic/gin"
)

type GetJobsResponse struct {
	Jobs []jobsService.Job `json:"jobs"`
}

//	@BasePath	/api/v1

//	@Summary	get jobs
//	@ID			get-jobs
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	GetJobsResponse	"ok"
//	@Router		/jobs [get]
func New(jobs *jobsService.Jobs) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobs, err := jobs.GetJobs()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"jobs": jobs})
	}
}
