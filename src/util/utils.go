package util

import (
	"log"
	"os"
	"strconv"
	"time"
)

// GetIntEnv accepts env variable name and default value.
// it it tries to get env variable, convert it to int and return parsed value otherwise default value is provided.
func GetIntEnv(envVariable string, defaultValue int) int {
	return GetEnv(envVariable, strconv.Atoi, defaultValue)
}

func GetDurationEnv(envVariable string, defaultValue time.Duration) time.Duration {
	return GetEnv(envVariable, time.ParseDuration, defaultValue)
}

func GetStringEnv(envVariable string, defaultValue string) string {
	return GetEnv(envVariable, func(val string) (string, error) {
		return val, nil
	}, defaultValue)
}

// GetEnv accepts env variable name, converter function and default value.
// it it tries to get env variable, convert it and return parsed value otherwise default value is provided.
func GetEnv[K interface{}](envVariable string, converter func(val string) (K, error), defaultValue K) K {
	if env := os.Getenv(envVariable); env != "" {
		if value, err := converter(env); err == nil {
			return value
		} else {
			log.Printf("Failed to parse value from environment variable %s: %s\n", envVariable, err)
		}
	}

	return defaultValue
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}
