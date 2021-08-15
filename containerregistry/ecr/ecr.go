package ecr

import "fmt"

func RegistryFormat(accountID string, region string) string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", accountID, region)
}
