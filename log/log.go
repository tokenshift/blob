package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	LEVEL_FATAL = "FATAL"
	LEVEL_ERROR = "ERROR"
	LEVEL_WARN  = "WARN "
	LEVEL_INFO  = "INFO "
	LEVEL_DEBUG = "DEBUG"
)

const LOG_TIMESTAMP_FORMAT = "2006-01-02T15:04:05Z07:00.000"

func Fatal(args...interface{}) {
	doLog(os.Stderr, LEVEL_FATAL, args...)
}

func Error(args...interface{}) {
	doLog(os.Stderr, LEVEL_ERROR, args...)
}

func Warn(args...interface{}) {
	doLog(os.Stderr, LEVEL_WARN, args...)
}

func Info(args...interface{}) {
	doLog(os.Stdout, LEVEL_INFO, args...)
}

func Debug(args...interface{}) {
	doLog(os.Stdout, LEVEL_DEBUG, args...)
}

func doLog(out io.Writer, level string, args...interface{}) {
	args = append([]interface{}{
		time.Now().Format(LOG_TIMESTAMP_FORMAT),
		level,
	}, args...)

	fmt.Fprintln(out, args...)
}
