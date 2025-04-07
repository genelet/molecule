package godbi

import (
	"io"
	"log/slog"
	"path/filepath"
)

type Slogger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

var _ Slogger = (*slog.Logger)(nil)

func DevelopLogger(w io.Writer, level ...slog.Level) *slog.Logger {
	levelvar := new(slog.LevelVar) // default Info. Debug, Info, Warn, Error
	if level != nil {
		levelvar.Set(level[0])
	} else { // reset default to be Debug
		levelvar.Set(slog.LevelDebug)
	}
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: levelvar, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}}))
}

func ProductLogger(w io.Writer, level ...slog.Level) *slog.Logger {
	levelvar := new(slog.LevelVar) // default Info
	if level != nil {
		levelvar.Set(level[0])
	}
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: levelvar}))
}
