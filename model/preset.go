package model

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/util"
)

type Preset struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Template string `json:"template"`
}

func GeneratePresetID() string {
	return fmt.Sprintf("pst-%s", util.RandomString(25))
}
