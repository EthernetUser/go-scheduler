package main

import (
	"context"
	"fmt"
	"go-scheduler/cmd/docs"
	"go-scheduler/internal/adapters"
	"go-scheduler/internal/config"
	"go-scheduler/internal/database/postgres"
	jobHandlers "go-scheduler/internal/handlers/jobs"
	jobService "go-scheduler/internal/jobs"
	gracefulshutdowner "go-scheduler/internal/pkg/gracefulShutdowner"
	slogLogger "go-scheduler/internal/pkg/logger"
	appWorker "go-scheduler/internal/worker"
	jobExecutor "go-scheduler/internal/worker/executor"
	robfigcron "go-scheduler/internal/worker/scheduler"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "go-scheduler/cmd/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	router := initRouterV1(logger, jobs)
	server := createHttpServer(cfg, router)
	gh := gracefulshutdowner.New(
		[]gracefulshutdowner.Instance{server, worker, db},
		cfg.GracefulSHutdowner.PerInstanceShutdownTimeout,
	)

	runApplication(logger, server, worker)

	logger.Debug("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeoutCause(
		context.Background(),
		cfg.GracefulSHutdowner.GlobalShutdownTimeout,
		gracefulshutdowner.GlobalTimeoutErr,
	)
	defer cancel()

	err = gh.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("Application has been stopped with errors", "errors", err)
		return
	}

	logger.Info("Application has been stopped successfully")
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
	logger *slog.Logger,
	jobs *jobService.Jobs,
) *http.ServeMux {
	router := http.NewServeMux()
	v1Preffix := "/api/v1"
	docs.SwaggerInfo.BasePath = v1Preffix
	router.Handle(
		"/swagger/",
		httpSwagger.WrapHandler,
	)

	slogAdapter := adapters.NewSlogAdapter(logger)
	handlers := jobHandlers.New(slogAdapter.WithGroup("service"), jobs)

	{
		router.Handle(fmt.Sprintf("POST %s/jobs", v1Preffix), handlers.Create())
		router.Handle(
			fmt.Sprintf(
				"DELETE %s/jobs/{%s}",
				v1Preffix,
				jobHandlers.JobIdPathParam,
			),
			handlers.Delete(),
		)
		router.Handle(fmt.Sprintf("GET %s/jobs", v1Preffix), handlers.GetAll())
		router.Handle(fmt.Sprintf("PUT %s/jobs", v1Preffix), handlers.Update())
	}
	return router
}
