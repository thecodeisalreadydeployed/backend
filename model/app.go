package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type App struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"project_id"`
	Name            string    `json:"name"`
	GitSource       GitSource `json:"git_source"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	BuildScript     string    `json:"build_script"`
	BuildCommand    string    `json:"build_command"`
	OutputDirectory string    `json:"output_directory"`
	InstallCommand  string    `json:"install_command"`
}

func GenerateAppID() string {
	return fmt.Sprintf("app_%s", util.RandomString(5))
}
