package create

import (
	jobsService "go-scheduler/internal/domain/jobs"

	"github.com/gin-gonic/gin"
)

type JobRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cron        string `json:"cron"`
	Url         string `json:"url"`
}

type JobResponse struct {
	Id          int    `json:"id"`
}

//	@BasePath	/api/v1

//	@Summary	Add a new job
//	@ID			create-job
//	@Accept		json
//	@Produce	json
//	@Param		request	body		JobRequest	true	"Some ID"
//	@Success	200		{object}	JobResponse		"ok"
//	@Router		/jobs [post]
func New(jobs *jobsService.Jobs) gin.HandlerFunc {
	return func(c *gin.Context) {
		var job JobRequest
		if err := c.BindJSON(&job); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		parsedJob := jobsService.JobForCreation{
			Name:        job.Name,
			Description: job.Description,
			Cron:        job.Cron,
			Url:         job.Url,
		}

		id, err := jobs.CreateJob(parsedJob)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"id": id})
	}
}
