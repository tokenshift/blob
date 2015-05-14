package env

import (
	"os"
	"strconv"
	"strings"
)

// Get the named environment variable.
func Get(name string) (string, bool) {
	for _, kv := range(os.Environ()) {
		split := strings.Index(kv, "=")
		if split > -1 && kv[0:split] == name {
			return kv[split+1:], true
		}
	}

	return "", false
}

// Get an environment variable as an integer.
func GetInt(name string) (int, bool, error) {
	value, ok := Get(name)
	if !ok {
		return -1, false, nil
	}

	i, err := strconv.ParseInt(value, 0, 0)
	if err != nil {
		return -1, true, err
	}

	return int(i), true, nil
}
