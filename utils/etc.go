package utils

import (
	"fmt"
	"os"
)

func Getenv(name string, defaultValue ...string) string {
	val := os.Getenv(name)

	if val == "" && len(defaultValue) > 0 {
		val = defaultValue[0]
	}

	return val
}

func InterfaceSliceToMapStringInterfaceSlice(i []interface{}) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	for _, v := range i {
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid value type (%v), it must be array of objects or just an object", v)
		}
		result = append(result, m)
	}

	return result, nil
}
