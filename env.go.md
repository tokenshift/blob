# Env

Utility functions for working with environment variables.

	<<#-->>
	package main

	import (
		"os"
		"strings"
	)

Environment variables are a set of key/value pairs of the form "key=value". In
order to retrieve specific variables by key, the input string is split on the
'=' sign.

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

GetEnv returns the requested variable and a bool indicating whether it was set;
GetEnvOr lets the consumer specify a default value to use if the variable was
not set.

	func GetEnv(key string) (string, bool) {
		for kv := range(Env()) {
			if kv.Key == key {
				return kv.Val, true
			}
		}

		return "", false
	}

	func GetEnvOr(key, alt string) string {
		if val, ok := GetEnv(key); ok {
			return val
		} else {
			return alt
		}
	}

By convention, lists of values in an environment variable are delimited using
the ':' character (e.g. $PATH). Colons can be escaped with a forward slash to
avoid using them as delimiters, if they need to be included in a value; for
example: HOSTS=localhost\:3000:example.com

	func GetEnvList(key string) ([]string, bool) {
		if val, ok := GetEnv(key); ok {
			return splitEnvList(val), true
		} else {
			return nil, false
		}
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
