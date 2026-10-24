package executor

import (
	"context"
	"fmt"
	"go-scheduler/internal/jobs"
	"net/http"
	"time"
)

type HTTPExecutor struct {
	client *http.Client
}

func NewHTTPExecutor(timeout time.Duration) *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{Timeout: timeout},
	}
}

func (e *HTTPExecutor) Execute(ctx context.Context, job jobs.Job) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.Url, nil)
	if err != nil {
		return fmt.Errorf("Error occurred on creating request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("Error occurred on processing request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}
