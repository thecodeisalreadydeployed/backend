package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type App struct {
	ID        string
	Name      string
	GitSource string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenerateAppID() string {
	return fmt.Sprintf("app_%s", util.RandomString(5))
}
