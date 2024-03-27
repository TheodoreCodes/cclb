package log

type Logger interface {
	Debug(msg string, args map[string]any)
	Info(msg string, args map[string]any)
	Warn(msg string, args map[string]any)
	Err(msg string, err error, args map[string]any)
}
