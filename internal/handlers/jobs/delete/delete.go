package delete

import (
	jobsService "go-scheduler/internal/domain/jobs"
	"strconv"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1

//	@Summary	delete a new job
//	@ID			delete-job
//	@Accept		json
//	@Produce	json
//	@Param		jobId	path		string	true	"Some ID"
//	@Success	200		{string}	string	"ok"
//	@Router		/jobs/{jobId} [delete]
func New(jobs *jobsService.Jobs) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobId := c.Param("jobId")
		if jobId == "" {
			c.JSON(400, gin.H{"error": "jobId is required"})
			return
		}

		parsedJobId, err := strconv.Atoi(jobId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = jobs.DeleteJob(parsedJobId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "success"})
	}
}
