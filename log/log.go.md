# Logging

	<<#-->>
	package log

	import (
		"fmt"
		"io"
		"os"
		"time"
	)

Simple logging functions used throughout the project.

	func log(out io.Writer, prefix string, args []interface{}) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(out, "%s %s ", timestamp, prefix)
		fmt.Fprintln(out, args...)
	}

	func Debug(args...interface{}) {
		log(os.Stdout, "DEBUG", args)
	}

	func Info(args...interface{}) {
		log(os.Stdout, "INFO ", args)
	}

	func Warn(args...interface{}) {
		log(os.Stderr, "WARN ", args)
	}

	func Error(args...interface{}) {
		log(os.Stderr, "ERROR", args)
	}

	func Fatal(args...interface{}) {
		log(os.Stderr, "FATAL", args)
	}
