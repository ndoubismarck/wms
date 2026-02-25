package types

type LogLevel string

const (
	LogLevelOff   LogLevel = "off"
	LogLevelInfo           = "infoLogo"
	LogLevelDebug          = "debug"
	LogLevelWarn           = "warn"
	LogLevelError          = "error"
)

type ILogger interface {
	Debug(v ...any)

	Debugf(format string, args ...interface{})

	Info(v ...any)

	Infof(format string, args ...interface{})

	Warn(v ...any)

	Warnf(format string, args ...interface{})

	Error(v ...any)

	Errorf(format string, args ...interface{})

	Fatal(v ...any)

	Fatalf(format string, args ...interface{})

	Panic(v ...any)

	Panicf(format string, args ...interface{})

	HTTP(v ...any)

	HTTPf(format string, args ...interface{})
}
