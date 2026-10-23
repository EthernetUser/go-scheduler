package main

import (
	"context"
	"go-scheduler/cmd/docs"
	"go-scheduler/internal/config"
	"go-scheduler/internal/database/postgres"
	jobHandlers "go-scheduler/internal/handlers/jobs"
	"go-scheduler/internal/jobs"
	"go-scheduler/internal/worker"
	"go-scheduler/internal/worker/executor"
	"go-scheduler/internal/worker/scheduler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.MustLoad()

	postgres, err := postgres.New(&cfg.Database)
	if err != nil {
		panic(err)
	}

	jobs := jobs.New(postgres)

	worker := initWorker(postgres)
	router := initRouterV1(jobs)

	server := &http.Server{
		Addr:         cfg.HttpServer.Addr,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			panic(err)
		}
	}()

	exitCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err := worker.Run(exitCtx); err != nil {
			panic(err)
		}
	}()

	<-exitCtx.Done()
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.HttpServer.ShutdownTimeout,
	)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error occurred: %v", err)
			return
		}
		log.Println("HttpServer has been stopped successfully")
	}()

	go func() {
		defer wg.Done()
		if err := worker.Shutdown(ctx); err != nil {
			log.Printf("Error occurred: %v", err)
			return
		}

		log.Println("Worker has been stopped successfully")
	}()

	wg.Wait()
	log.Println("Application has beeb stopped")
}

func initWorker(postgres *postgres.Postgres) *worker.Worker {
	cron := scheduler.NewRobfigCronScheduler()
	httpExecutor := executor.NewHTTPExecutor(30 * time.Second)

	return worker.New(postgres, httpExecutor, 2*time.Minute, cron)
}

func initRouterV1(jobs *jobs.Jobs) *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/jobs", jobHandlers.CreateCreateHandler(jobs))
		v1.DELETE("/jobs/:jobId", jobHandlers.CreateDeleteHandler(jobs))
		v1.GET("/jobs", jobHandlers.CreateGetHandler(jobs))
	}
	return router
}
