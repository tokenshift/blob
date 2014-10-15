# Logging

Utility functions for logging. These need to be refactored to allow logging
from anywhere in the app while tracking a request ID.

	<<#-->>

	package main

	import (
		"fmt"
	)

The `HasID` context provides an ID for tracking a specific request throughout
the app.

	type HasID interface {
		ID() int
	}

All logs are written to STDOUT, with a log level included in the message.

	func log(ctx HasID, level string, args...interface{}) {
		args = append([]interface{}{fmt.Sprintf("%8x %s:", ctx.ID(), level)}, args...)
		fmt.Println(args...)
	}

	func LogError(ctx HasID, args...interface{}) {
		log(ctx, "ERROR", args...)
	}

	func LogInfo(ctx HasID, args...interface{}) {
		log(ctx, "INFO", args...)
	}

	func LogWarn(ctx HasID, args...interface{}) {
		log(ctx, "WARN", args...)
	}
