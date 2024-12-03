package main

import (
	"go-scheduler/cmd/docs"
	"go-scheduler/internal/config"
	"go-scheduler/internal/database/postgres"
	"go-scheduler/internal/domain/jobs"
	"go-scheduler/internal/handlers/jobs/create"
	"go-scheduler/internal/handlers/jobs/delete"
	"go-scheduler/internal/handlers/jobs/get"
	"go-scheduler/internal/pkg/cronScheduler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.MustLoad()

	postgres, err := postgres.New(&cfg.Database)
	if err != nil {
		panic(err)
	}
	
	scheduler := cronScheduler.New()

	jobs := jobs.New(postgres, scheduler)

	router := initRouterV1(jobs)

	server := &http.Server{
		Addr: cfg.HttpServer.Addr,
		ReadTimeout: cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		Handler: router,
	}

	server.ListenAndServe()
}


func initRouterV1(jobs *jobs.Jobs) *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/jobs", create.New(jobs))
		v1.DELETE("/jobs/:jobId", delete.New(jobs))
		v1.GET("/jobs", get.New(jobs))
	}
	return router
}