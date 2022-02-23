package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GenerateProjectID() string {
	return fmt.Sprintf("prj-%s", util.RandomString(7))
}
