package main

import (
	"os"
	"strings"
)

type KeyVal struct {
	Key, Val string
}

func Env() (<- chan KeyVal) {
	out := make(chan KeyVal)

	go func() {
		defer close(out)
		for _, entry := range(os.Environ()) {
			split := strings.SplitN(entry, "=", 2)
			if len(split) == 2 {
				out <- KeyVal{split[0],split[1]}
			}
		}
	}()

	return out
}

func GetEnv(key string) (string, bool) {
	for kv := range(Env()) {
		if kv.Key == key {
			return kv.Val, true
		}
	}

	return "", false
}

func splitEnvList(val string) []string {
	list := make([]string, 0)

	var i, start int
	for i, start = 0, 0; i < len(val); i += 1 {
		c := val[i]

		// Escape character
		if c == '\\' && len(val) > i+1 {
			val = val[:i] + val[i+1:]
			continue
		}

		if c == ':' {
			list = append(list, val[start:i])
			start = i+1
		}
	}

	if i > start {
		list = append(list, val[start:i])
	}

	return list
}

func GetEnvList(key string) ([]string, bool) {
	if val, ok := GetEnv(key); ok {
		return splitEnvList(val), true
	} else {
		return nil, false
	}
}

func GetEnvOr(key, alt string) string {
	if val, ok := GetEnv(key); ok {
		return val
	} else {
		return alt
	}
}
