package worker

type Cron interface {
	Start()
	Stop()
	AddFunc(cron string, cmd func()) (int, error)
	Remove(int)
}
