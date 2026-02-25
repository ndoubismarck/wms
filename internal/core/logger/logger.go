package logger

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"server/internal/core/shared/types"
	"strings"

	log "github.com/sirupsen/logrus"
)

//ToDo: revisit and fix loglevel, log config and logger hook

type logger struct {
	debug *log.Logger
	warn  *log.Logger
	info  *log.Logger
	http  *log.Logger
	error *log.Logger
	fatal *log.Logger
	panic *log.Logger
}

var logLevel = types.LogLevel(types.LogLevelDebug)

func New(ctx types.IContext) types.ILogger {
	env, err := ctx.DotEnv().Get()
	if err == nil {
		logLevel = types.LogLevel(env.Log.Level)
	}
	debugLog := log.New()
	debugLog.SetLevel(log.DebugLevel)
	debugLog.SetOutput(io.Discard)
	debugLog.AddHook(newHook("debug"))

	infoLog := log.New()
	infoLog.SetLevel(log.DebugLevel)
	infoLog.SetOutput(io.Discard)
	infoLog.AddHook(newHook("info"))

	warnLog := log.New()
	warnLog.SetLevel(log.DebugLevel)
	warnLog.SetOutput(io.Discard)
	warnLog.AddHook(newHook("warn"))

	errorLog := log.New()
	errorLog.SetLevel(log.DebugLevel)
	errorLog.SetOutput(io.Discard)
	errorLog.AddHook(newHook("error"))

	fatalLog := log.New()
	fatalLog.SetLevel(log.DebugLevel)
	fatalLog.SetOutput(io.Discard)
	fatalLog.AddHook(newHook("fatal"))

	httpLog := log.New()
	httpLog.SetLevel(log.DebugLevel)
	httpLog.SetOutput(io.Discard)
	httpLog.AddHook(newHook("http"))

	return &logger{
		info:  infoLog,
		warn:  warnLog,
		debug: debugLog,
		fatal: fatalLog,
		error: errorLog,
		http:  httpLog,
	}
}

func (l *logger) Debug(v ...any) {
	v = append(v, l.source())
	l.debug.Debug(v...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	format = format + l.source()
	l.debug.Debugf(format, args...)
}

func (l *logger) Info(v ...any) {
	v = append(v, l.source())
	l.info.Info(v...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	format = format + l.source()
	l.info.Infof(format, args...)
}

func (l *logger) Warn(v ...any) {
	v = append(v, l.source())
	l.warn.Warn(v...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	format = format + l.source()
	l.warn.Warnf(format, args...)
}

func (l *logger) Error(v ...any) {
	v = append(v, l.source())
	l.error.Error(v...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	format = format + l.source()
	l.error.Errorf(format, args...)
}

func (l *logger) Fatal(v ...any) {
	v = append(v, l.source())
	l.fatal.Fatal(v...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	format = format + l.source()
	l.fatal.Fatalf(format, args...)
}

func (l *logger) Panic(v ...any) {
	v = append(v, l.source())
	l.panic.Panic(v...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	format = format + l.source()
	l.panic.Panicf(format, args...)
}

func (l *logger) HTTP(v ...any) {
	l.http.Debug(v...)
}

func (l *logger) HTTPf(format string, args ...interface{}) {
	l.http.Debugf(format, args...)
}

func (l *logger) source() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		if strings.Contains(file, "/cmd/") {
			val := strings.Split(file, "/cmd/")
			if len(val) >= 2 && strings.HasSuffix(val[1], ".go") {
				file = filepath.Join("cmd", val[1])
			}
		}
		if strings.Contains(file, "/pkg/") {
			val := strings.Split(file, "/pkg/")
			if len(val) >= 2 && strings.HasSuffix(val[1], ".go") {
				file = filepath.Join("pkg", val[1])
			}
		}
		if strings.Contains(file, "/internal/") {
			val := strings.Split(file, "/internal/")
			if len(val) >= 2 && strings.HasSuffix(val[1], ".go") {
				file = filepath.Join("internal", val[1])
			}
		}
		return fmt.Sprintf(" | %s:%d", file, line)
	}
	return ""
}
