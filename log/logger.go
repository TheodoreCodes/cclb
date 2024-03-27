package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type DefaultLogger struct {
}

func (rv *DefaultLogger) addArgs(ev *zerolog.Event, args map[string]any) {
	for k, v := range args {
		switch val := v.(type) {
		case string:
			ev.Str(k, val)
			break
		case int:
			ev.Int(k, val)
			break
		case float64:
			ev.Float64(k, val)
			break
		case bool:
			ev.Bool(k, val)
		}
	}
}

func (rv *DefaultLogger) Debug(msg string, args map[string]any) {
	ev := log.Debug()
	rv.addArgs(ev, args)
	ev.Msg(msg)
}

func (rv *DefaultLogger) Info(msg string, args map[string]any) {
	ev := log.Info()

	rv.addArgs(ev, args)
	ev.Msg(msg)
}

func (rv *DefaultLogger) Warn(msg string, args map[string]any) {
	ev := log.Warn()

	rv.addArgs(ev, args)
	ev.Msg(msg)
}

func (rv *DefaultLogger) Err(msg string, err error, args map[string]any) {
	ev := log.Error()
	ev.Err(err)
	rv.addArgs(ev, args)
	ev.Msg(msg)
}

func NewDefaultLogger(logLevel string) Logger {
	logOutput := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(logOutput)

	var l zerolog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		l = zerolog.DebugLevel
		break
	case "info":
		l = zerolog.InfoLevel
		break
	case "warn":
		l = zerolog.WarnLevel
		break
	case "error":
		l = zerolog.ErrorLevel
		break
	default:
		panic("illegal logLevel value")
	}

	zerolog.SetGlobalLevel(l)

	return &DefaultLogger{}
}
