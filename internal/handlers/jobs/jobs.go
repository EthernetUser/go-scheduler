package handlers

type Logger interface {
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	WithGroup(name string) Logger
}

type handlers struct {
	logger Logger
}

func New(logger Logger) *handlers {
	return &handlers{
		logger: logger,
	}
}
