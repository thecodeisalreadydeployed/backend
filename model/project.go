package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GenerateProjectID() string {
	return fmt.Sprintf("prj_%s", util.RandomString(5))
}
