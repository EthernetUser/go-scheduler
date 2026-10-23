package main

import (
	"context"
	"go-scheduler/cmd/docs"
	"go-scheduler/internal/config"
	"go-scheduler/internal/database/postgres"
	jobHandlers "go-scheduler/internal/handlers/jobs"
	jobService "go-scheduler/internal/jobs"
	appWorker "go-scheduler/internal/worker"
	jobExecutor "go-scheduler/internal/worker/executor"
	robfigcron "go-scheduler/internal/worker/scheduler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.MustLoad()

	db, err := postgres.New(&cfg.Database)
	if err != nil {
		log.Fatalf("db connection faild: %v", err)
	}

	jobs := jobService.New(db)
	worker := initWorker(db, cfg)
	router := initRouterV1(jobs)
	server := createHttpServer(cfg, router)

	runApplication(server, worker)

	log.Println("Shutting down server...")

	shutdownApplication(cfg, server, worker, db)

	log.Println("Application has been stopped")
}

func runApplication(
	server *http.Server,
	worker *appWorker.Worker,
) {
	exitCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 2)
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	go func() {
		if err := worker.Run(exitCtx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-exitCtx.Done():
	case err := <-errCh:
		log.Printf("Component failed: %v", err)
	}
}

func shutdownApplication(
	cfg *config.Config,
	server *http.Server,
	worker *appWorker.Worker,
	db *postgres.Postgres,
) {
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

	if err := db.Close(); err != nil {
		log.Printf("db shutting down is failed, err: %v", err)
		return
	}
	log.Println("DB connection has been closed successfully")
}

func initWorker(
	db *postgres.Postgres,
	cfg *config.Config,
) *appWorker.Worker {
	cron := robfigcron.New()
	httpExecutor := jobExecutor.NewHTTPExecutor(cfg.Worker.JobTimeout)

	return appWorker.New(
		db,
		httpExecutor,
		cfg.Worker.ReconclineInterval,
		cron,
		cfg.Worker.JobTimeout,
	)
}

func createHttpServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.HttpServer.Addr,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		Handler:      router,
	}
}

func initRouterV1(jobs *jobService.Jobs) *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/jobs", jobHandlers.Create(jobs))
		v1.DELETE("/jobs/:jobId", jobHandlers.Delete(jobs))
		v1.GET("/jobs", jobHandlers.Get(jobs))
	}
	return router
}
