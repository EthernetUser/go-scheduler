package main

import (
	"context"
	"go-scheduler/cmd/docs"
	"go-scheduler/internal/adapters"
	"go-scheduler/internal/config"
	"go-scheduler/internal/database/postgres"
	jobHandlers "go-scheduler/internal/handlers/jobs"
	jobService "go-scheduler/internal/jobs"
	slogLogger "go-scheduler/internal/pkg/logger"
	appWorker "go-scheduler/internal/worker"
	jobExecutor "go-scheduler/internal/worker/executor"
	robfigcron "go-scheduler/internal/worker/scheduler"
	"io"
	"log/slog"
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
	logger := slogLogger.New(cfg.Env)
	db, err := postgres.New(&cfg.Database)

	if err != nil {
		logger.Error("DB connection failed", "error", err)
		os.Exit(1)
	}

	jobs := jobService.New(db)
	worker := initWorker(db, cfg, logger)
	router := initRouterV1(cfg.Env, logger, jobs)
	server := createHttpServer(cfg, router)

	runApplication(logger, server, worker)

	logger.Debug("Shutting down server...")

	shutdownApplication(cfg, logger, server, worker, db)
	logger.Info("Application has been stopped")
}

func runApplication(
	logger *slog.Logger,
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
		logger.Error("Running application failed", "error", err)
	}
}

func shutdownApplication(
	cfg *config.Config,
	logger *slog.Logger,
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
			logger.Error("Error occurred", "error", err)
			return
		}
		logger.Debug("HttpServer has been stopped successfully")
	}()

	go func() {
		defer wg.Done()
		if err := worker.Shutdown(ctx); err != nil {
			logger.Error("Error occurred", "error", err)
			return
		}

		logger.Debug("Worker has been stopped successfully")
	}()

	wg.Wait()

	if err := db.Close(); err != nil {
		logger.Error("DB shutting down is failed", "error", err)
		return
	}
	logger.Debug("DB connection has been closed successfully")
}

func initWorker(
	db *postgres.Postgres,
	cfg *config.Config,
	logger *slog.Logger,
) *appWorker.Worker {
	cron := robfigcron.New()
	httpExecutor := jobExecutor.NewHTTPExecutor(cfg.Worker.JobTimeout)

	return appWorker.New(
		db,
		logger.WithGroup("worker"),
		httpExecutor,
		cfg.Worker.ReconclineInterval,
		cron,
		cfg.Worker.JobTimeout,
	)
}

func createHttpServer(
	cfg *config.Config,
	router http.Handler,
) *http.Server {
	return &http.Server{
		Addr:         cfg.HttpServer.Addr,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		Handler:      router,
	}
}

func initRouterV1(
	env string,
	logger *slog.Logger,
	jobs *jobService.Jobs,
) *gin.Engine {
	gin.DefaultWriter = io.Discard
	switch env {
	case "local":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	slogAdapter := adapters.NewSlogAdapter(logger)
	handlers := jobHandlers.New(slogAdapter.WithGroup("service"))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/jobs", handlers.Create(jobs))
		v1.DELETE("/jobs/:jobId", handlers.Delete(jobs))
		v1.GET("/jobs", handlers.Get(jobs))
		v1.PUT("/jobs", handlers.Update(jobs))
	}
	return router
}
