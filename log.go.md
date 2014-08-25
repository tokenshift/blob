# Blob - Logging

	<<#-->>
	package main

	import "fmt"
	import "os"
	import "time"

Log output is routed through a channel to ensure messages from multiple threads
are not interleaved, the written to STDOUT or STDERR depending on the severity.

	type logSeverity int
	const (
		debug logSeverity = iota
		info
		warn
		error
		fatal
	)

	type logMessage struct {
		timestamp time.Time
		severity logSeverity
		message string
	}

	var logChannel = make(chan logMessage, 100)

	func log(severity logSeverity, message string) {
		logChannel <- logMessage {
			timestamp: time.Now(),
			severity: severity,
			message: message,
		}
	}

	func logDebug(msgs...interface{}) {
		log(debug, fmt.Sprint(msgs...))
	}

	func logDebugf(format string, args...interface{}) {
		log(debug, fmt.Sprintf(format, args...))
	}

	func logInfo(msgs...interface{}) {
		log(info, fmt.Sprint(msgs...))
	}

	func logInfof(format string, args...interface{}) {
		log(info, fmt.Sprintf(format, args...))
	}

	func logWarn(msgs...interface{}) {
		log(warn, fmt.Sprint(msgs...))
	}

	func logWarnf(format string, args...interface{}) {
		log(warn, fmt.Sprintf(format, args...))
	}

	func logError(msgs...interface{}) {
		log(error, fmt.Sprint(msgs...))
	}

	func logErrorf(format string, args...interface{}) {
		log(error, fmt.Sprintf(format, args...))
	}

	func logFatal(msgs...interface{}) {
		log(fatal, fmt.Sprint(msgs...))
	}

	func logFatalf(format string, args...interface{}) {
		log(fatal, fmt.Sprintf(format, args...))
	}

	func runLogger() {
		for msg := range(logChannel) {
			timestamp := msg.timestamp.Format(time.RFC3339)
			switch msg.severity {
			case debug:
				fmt.Fprintln(os.Stdout, timestamp, "DEBUG:", msg.message)
			case info:
				fmt.Fprintln(os.Stdout, timestamp, " INFO:", msg.message)
			case warn:
				fmt.Fprintln(os.Stderr, timestamp, " WARN:", msg.message)
			case error:
				fmt.Fprintln(os.Stderr, timestamp, "ERROR:", msg.message)
			case fatal:
				fmt.Fprintln(os.Stderr, timestamp, "FATAL:", msg.message)
			}
		}
	}
