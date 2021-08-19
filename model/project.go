package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Project struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenerateProjectID() string {
	return fmt.Sprintf("prj_%s", util.RandomString(5))
}
