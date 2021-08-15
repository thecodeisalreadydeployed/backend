package gcr

import (
	"errors"
	"fmt"
	"strings"
)

func RegistryFormat(hostname string, projectID string) (string, error) {
	if !strings.Contains(hostname, "gcr.io") {
		return "", errors.New("Invalid hostname for Google Container Registry.")
	}

	return fmt.Sprintf("%s/%s", hostname, projectID), nil
}
