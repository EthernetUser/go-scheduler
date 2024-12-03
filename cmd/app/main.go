package main

import (
	"context"
	"go-scheduler/cmd/docs"
	"go-scheduler/internal/config"
	"go-scheduler/internal/database/postgres"
	"go-scheduler/internal/domain/jobs"
	"go-scheduler/internal/handlers/jobs/create"
	"go-scheduler/internal/handlers/jobs/delete"
	"go-scheduler/internal/handlers/jobs/get"
	"go-scheduler/internal/pkg/cronScheduler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	go func () {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<- quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped")
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