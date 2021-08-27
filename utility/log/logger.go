package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	defaultLevel = 2
	splitter     = "thrive-challenges"
)

type Level string

const (
	LevelUnknown = Level("unknown")
	LevelDebug   = Level("debug")
	LevelInfo    = Level("info")
	LevelWarn    = Level("warn")
	LevelError   = Level("error")
	LevelPanic   = Level("panic")
)

var levels = map[string]Level{
	LevelDebug.String(): LevelDebug,
	LevelInfo.String():  LevelInfo,
	LevelWarn.String():  LevelWarn,
	LevelError.String(): LevelError,
	LevelPanic.String(): LevelPanic,
}

func ParseLevel(level string) Level {
	if l, ok := levels[level]; ok {
		return l
	}
	return LevelInfo
}

func (l Level) String() string {
	return string(l)
}

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Panic(format string, args ...interface{})

	SetLevel(level Level)
}

type logger struct {
	l          *logrus.Logger
	stackLevel int
}

func New(level Level) *logger {
	return newWithStackLevel(level, defaultLevel)
}

// create a new logrus logger and also set the level
func newWithStackLevel(level Level, stackLevel int) *logger {
	l := &logger{
		l:          logrus.New(),
		stackLevel: stackLevel,
	}
	l.SetLevel(level)
	return l
}

func (l *logger) log(level Level, format string, args ...interface{}) {
	var msg string
	if len(args) == 0 {
		msg = format
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	entry := l.l.WithField("msg", msg)

	// Append caller line number if possible
	_, file, line, ok := runtime.Caller(l.stackLevel)
	parts := strings.Split(file, splitter)
	if ok && len(parts) > 1 {
		entry = entry.WithField("location", fmt.Sprintf("%s:%d", parts[1], line))
	}

	// this is used to add log level to printed logs
	switch level {
	case LevelDebug:
		entry.Debug("")
	case LevelInfo:
		entry.Info("")
	case LevelWarn:
		entry.Warn("")
	case LevelError:
		entry.Error("")
	case LevelPanic:
		entry.Panic("")
	default:
		entry.Debug("")
	}

}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

func (l *logger) Warning(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

func (l *logger) Panic(format string, args ...interface{}) {
	l.log(LevelPanic, format, args...)
}

func (l *logger) SetLevel(logLevel Level) {
	level, err := logrus.ParseLevel(logLevel.String())
	if err == nil {
		l.l.SetLevel(level)
		return
	}
	l.Panic(err.Error())
}

func (l *logger) SetWriter(w io.Writer) {
	l.l.Out = w
}
