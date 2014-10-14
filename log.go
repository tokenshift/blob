package main

import (
	"fmt"
)

func log(level string, args...interface{}) {
	args = append([]interface{}{fmt.Sprintf("%s:", level)}, args...)
	fmt.Println(args...)
}

func logError(args...interface{}) {
	log("ERROR", args...)
}

func logInfo(args...interface{}) {
	log("INFO", args...)
}

func logWarn(args...interface{}) {
	log("WARN", args...)
}
