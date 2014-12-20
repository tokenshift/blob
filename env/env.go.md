# Env

Utility functions for working with environment variables.

	<<#-->>
	package env

	import (
		"encoding/base32"
		"encoding/base64"
		"encoding/hex"
		"os"
		"strconv"
		"strings"
	)

Environment variables are a set of key/value pairs of the form "key=value". In
order to retrieve specific variables by key, the input string is split on the
'=' sign.

	type KeyVal struct {
		Key, Val string
	}


	func All() (<- chan KeyVal) {
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

Get returns the requested variable and a bool indicating whether it was set;
GetOr lets the consumer specify a default value to use if the variable was
not set.

	func Get(key string) (string, bool) {
		for kv := range(All()) {
			if kv.Key == key {
				return kv.Val, true
			}
		}

		return "", false
	}

	func GetOr(key, alt string) string {
		if val, ok := Get(key); ok {
			return val
		} else {
			return alt
		}
	}

	func GetInt(key string) (int, bool, error) {
		if val, ok := Get(key); ok {
			i, err := strconv.ParseInt(val, 0, 0)
			return int(i), true, err
		} else {
			return -1, false, nil
		}
	}

By convention, lists of values in an environment variable are delimited using
the ':' character (e.g. $PATH). Colons can be escaped with a forward slash to
avoid using them as delimiters, if they need to be included in a value; for
example: HOSTS=localhost\:3000:example.com

	func GetList(key string) ([]string, bool) {
		if val, ok := Get(key); ok {
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

Binary data can also be provided as environment variables using a specific
encoding. Currently, base 16 (hexadecimal), 32, and 64 are supported. These
functions return false if the specified value was not found, or true and an
error if the value was found but could not be decoded.

	func Get16(key string) ([]byte, bool, error) {
		input, ok := Get(key)
		if !ok {
			return nil, false, nil
		}

		decoded, err := hex.DecodeString(input)
		return decoded, true, err
	}

	func Get32(key string) ([]byte, bool, error) {
		input, ok := Get(key)
		if !ok {
			return nil, false, nil
		}

		decoded, err := base32.StdEncoding.DecodeString(input)
		return decoded, true, err
	}

	func Get64(key string) ([]byte, bool, error) {
		input, ok := Get(key)
		if !ok {
			return nil, false, nil
		}

		decoded, err := base64.StdEncoding.DecodeString(input)
		return decoded, true, err
	}
