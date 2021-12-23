package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type App struct {
	ID                 string             `json:"id"`
	ProjectID          string             `json:"projectID"`
	Name               string             `json:"name"`
	GitSource          GitSource          `json:"gitSource"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	BuildConfiguration BuildConfiguration `json:"buildConfiguration"`
	Observable         bool               `json:"observable"`
}

func GenerateAppID() string {
	return fmt.Sprintf("app_%s", util.RandomString(25))
}
