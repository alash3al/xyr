package utils

import "os"

func Getenv(name string, defaultValue ...string) string {
	val := os.Getenv(name)

	if val == "" && len(defaultValue) > 0 {
		val = defaultValue[0]
	}

	return val
}
