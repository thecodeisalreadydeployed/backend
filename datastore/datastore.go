package datastore

import (
	"fmt"
)

func GetEvent(id string) string {
	return fmt.Sprintf("Dummy event %s.", id)
}
