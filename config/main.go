package config

import (
	"fmt"
	"os"
)

func GetProjectId() (string, error) {
	projectId := os.Getenv("PROJECT_ID")

	if projectId == "" {
		return "", fmt.Errorf("project ID is empty")
	}

	return projectId, nil
}

func GetCredentials() (string, error) {
	serviceAccount := os.Getenv("CREDENTIALS")

	if serviceAccount == "" {
		return "", fmt.Errorf("credentials is empty")
	}

	return serviceAccount, nil
}
