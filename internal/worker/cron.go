package worker

import "context"

type Cron interface {
	Start()
	Stop() context.Context
	AddFunc(cron string, cmd func()) (int, error)
	Remove(int)
}
