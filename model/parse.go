package model

import (
	"encoding/base64"
	"encoding/json"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func GetCreatorString(c Actor) string {
	creator, err := json.Marshal(c)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return cast.ToString(creator)
}

func GetGitSourceString(gs GitSource) string {
	gitSource, err := json.Marshal(gs)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return cast.ToString(gitSource)
}

func GetBuildConfigurationString(bc BuildConfiguration) string {
	buildConfiguration, err := json.Marshal(bc)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return base64.StdEncoding.EncodeToString(buildConfiguration)
}
