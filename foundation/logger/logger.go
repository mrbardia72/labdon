// Package logger provides a minimal structured leveled logger.
// It is deliberately small and interface-based so it can be swapped for
// zap (or anything else) without touching call sites.
package logger

import (
	"context"
	"io"
	"log"
	"os"
)

// Level identifies the severity of a log message.
type Level int

// Supported levels, ordered from most to least severe.
const (
	ErrorLevel Level = iota + 1
	WarnLevel
	InfoLevel
	DebugLevel
)

// String implements fmt.Stringer.
func (l Level) String() string {
	switch l {
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Logger is the interface the rest of the code depends on.
type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, msg string, keysAndValues ...any)
	Error(ctx context.Context, msg string, keysAndValues ...any)
}

type stdLogger struct {
	out   *log.Logger
	level Level
}

// NewDev returns a human-readable logger writing JSON-ish key pairs to w
// (stderr when w is nil) at the given level.
func NewDev(w io.Writer, level Level) Logger {
	if w == nil {
		w = os.Stderr
	}
	return &stdLogger{out: log.New(w, "", log.LstdFlags|log.Lmsgprefix), level: level}
}

func (l *stdLogger) log(level Level, ctx context.Context, msg string, kv ...any) {
	if level > l.level {
		return
	}
	l.out.Printf("%s %s %v", level, msg, kv)
}

func (l *stdLogger) Debug(ctx context.Context, msg string, kv ...any) {
	l.log(DebugLevel, ctx, msg, kv...)
}

func (l *stdLogger) Info(ctx context.Context, msg string, kv ...any) {
	l.log(InfoLevel, ctx, msg, kv...)
}

func (l *stdLogger) Warn(ctx context.Context, msg string, kv ...any) {
	l.log(WarnLevel, ctx, msg, kv...)
}

func (l *stdLogger) Error(ctx context.Context, msg string, kv ...any) {
	l.log(ErrorLevel, ctx, msg, kv...)
}
