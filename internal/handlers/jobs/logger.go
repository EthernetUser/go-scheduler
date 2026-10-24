package handlers

type Logger interface {
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
