package handlers

import (
	jobsService "go-scheduler/internal/jobs"
	"strconv"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1

// @Summary	delete a new job
// @ID			delete-job
// @Accept		json
// @Produce	json
// @Param		jobId	path		string	true	"Some ID"
// @Success	200		{string}	string	"ok"
// @Router		/jobs/{jobId} [delete]
func (h *handlers) Delete(jobs *jobsService.Jobs) gin.HandlerFunc {
	logger := h.logger.WithGroup("handler.delete")
	return func(ctx *gin.Context) {
		jobId := ctx.Param("jobId")
		if jobId == "" {
			ctx.JSON(400, gin.H{"error": "jobId is required"})
			return
		}

		parsedJobId, err := strconv.Atoi(jobId)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = jobs.DeleteJob(ctx, parsedJobId)
		if err != nil {
			logger.Error("Error occurred", "error", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "success"})
	}
}
