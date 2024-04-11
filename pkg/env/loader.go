package env

import (
	"fmt"
	"os"
	"strings"
)

func GetConfigPath(suffix string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	appEnv := getEnv("APP_ENV", "dev")
	if strings.HasSuffix(dir, "config") {
		return fmt.Sprintf("%s/config.%s.yaml", dir, appEnv), nil
	} else if strings.HasSuffix(dir, suffix) {
		return fmt.Sprintf("%s/config.%s.yaml", dir, appEnv), nil
	} else if strings.HasSuffix(dir, "interview") {
		return fmt.Sprintf("%s/cmd/%s/config/config.%s.yaml", dir, suffix, appEnv), nil
	} else {
		splitted := strings.Split(dir, suffix)
		return fmt.Sprintf("%s/%s/config/config.%s.yaml", splitted[0], suffix, appEnv), nil
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
