package main

import (
	"fmt"
)

type HasID interface {
	ID() int
}

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
