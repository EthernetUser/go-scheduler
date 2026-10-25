package gracefulshutdowner

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var PerInstanceTimeoutErr = errors.New("perInstanceTimeout execeed")
var GlobalTimeoutErr = errors.New("globalTimeout execeed")

type Instance interface {
	Shutdown(context.Context) error
}

type gracefulshutdowner struct {
	instances          []Instance
	timeoutPerInstance time.Duration
}

func New(
	instances []Instance,
	timeoutPerInstance time.Duration,
) *gracefulshutdowner {
	return &gracefulshutdowner{
		instances:          instances,
		timeoutPerInstance: timeoutPerInstance,
	}
}

func (g *gracefulshutdowner) Shutdown(ctx context.Context) error {
	var errs = make([]error, 0, len(g.instances))

	for _, instance := range g.instances {
		err := g.shutdownOnce(ctx, instance)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (g *gracefulshutdowner) shutdownOnce(
	ctx context.Context,
	instance Instance,
) error {
	timoutPreInstanceCtx, cancel := context.WithTimeoutCause(
		ctx,
		g.timeoutPerInstance,
		PerInstanceTimeoutErr,
	)
	defer cancel()

	if err := instance.Shutdown(timoutPreInstanceCtx); err != nil {
		return fmt.Errorf("%T: %w", instance, err)
	}

	return nil
}
