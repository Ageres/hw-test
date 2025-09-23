package utils

import (
	"log"
	"os"
)

func GetEnvOrDefault(envName string, defaultVal string) string {
	envValue, isSet := os.LookupEnv(envName)
	if isSet {
		return envValue
	}
	log.Println("Environment variable is not provided, using default value",
		map[string]any{
			"variableName": envName,
			"defaultValue": defaultVal,
		},
	)
	return defaultVal
}
