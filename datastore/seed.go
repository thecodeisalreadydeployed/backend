package datastore

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/preset"
	"math/rand"

	faker "github.com/bxcodec/faker/v3"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func seedPreset() {
	if seedExists("presets") {
		return
	}

	var data []datamodel.Preset
	data_ := []model.Preset{
		{
			ID:       "pst-flaskframeworkpresetxxxxx",
			Name:     "Flask Default Preset",
			Template: preset.Text(preset.FrameworkFlask),
		},
		{
			ID:       "pst-springframeworkpresetxxxx",
			Name:     "Spring Default Preset",
			Template: preset.Text(preset.FrameworkSpring),
		},
		{
			ID:       "pst-nestjsframeworkpresetxxxx",
			Name:     "NestJS Default Preset",
			Template: preset.Text(preset.FrameworkNestJS),
		},
		{
			ID:       "pst-simplepresetxxxxxxxxxxxxx",
			Name:     "Simple Preset",
			Template: preset.Text(preset.NoFramework),
		},
	}

	for _, datum_ := range data_ {
		datum := *datamodel.NewPresetFromModel(&datum_)
		data = append(data, datum)
	}

	if err := GetDB().Omit("Deployment").Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed apps.")
	}
}

func seed() {
	seedProjects(5)
	seedApps(15)
	seedDeployments(40)
	seedEvents(60)
}

func seedExists(name string) bool {
	var existing int64
	err := GetDB().Table(name).Count(&existing).Error
	if err != nil {
		zap.L().Error(err.Error())
		return false
	} else if existing > 0 {
		zap.L().Info(fmt.Sprintf("Table '%s' already seeded.", name))
		return true
	}
	return false
}

func seedProjects(size int) {
	if seedExists("projects") {
		return
	}

	var data []datamodel.Project
	for i := 0; i < size; i++ {
		var datum datamodel.Project
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = withPrefix(datum.ID, "prj")
		data = append(data, datum)
	}
	if err := GetDB().Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed projects.")
	}

}

func seedApps(size int) {
	if seedExists("apps") {
		return
	}

	var keys []string
	err := GetDB().Table("projects").Select("ID").Scan(&keys).Error
	if err != nil {
		zap.L().Error(err.Error())
	}

	var data []datamodel.App
	for i := 0; i < size; i++ {
		var datum datamodel.App
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = withPrefix(datum.ID, "app")
		datum.ProjectID = getForeignKey(keys)
		datum.GitSource = getGitSource()
		datum.BuildConfiguration = getBuildConfiguration()
		datum.Observable = false

		data = append(data, datum)
	}
	if err := GetDB().Omit("Project").Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed apps.")
	}
}

func seedDeployments(size int) {
	if seedExists("deployments") {
		return
	}

	var keys []string
	err := GetDB().Table("apps").Select("ID").Scan(&keys).Error
	if err != nil {
		zap.L().Error(err.Error())
	}

	var data []datamodel.Deployment
	for i := 0; i < size; i++ {
		var datum datamodel.Deployment
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = withPrefix(datum.ID, "dpl")
		datum.AppID = getForeignKey(keys)
		datum.GitSource = getGitSource()
		datum.Creator = getCreator()
		datum.BuildConfiguration = getBuildConfiguration()
		datum.State = model.DeploymentState(getState())

		data = append(data, datum)
	}
	if err := GetDB().Omit("App").Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed deployments.")
	}

}

func seedEvents(size int) {
	if seedExists("events") {
		return
	}

	var keys []string
	err := GetDB().Table("deployments").Select("ID").Scan(&keys).Error
	if err != nil {
		zap.L().Error(err.Error())
	}

	var data []datamodel.Event
	for i := 0; i < size; i++ {
		var datum datamodel.Event
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = model.GenerateEventID(datum.ExportedAt)
		datum.DeploymentID = getForeignKey(keys)
		datum.Type = model.EventType(getType())

		data = append(data, datum)
	}
	if err := GetDB().Omit("Deployment").Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed events.")
	}
}

func getForeignKey(keys []string) string {
	return keys[rand.Intn(len(keys))]
}

func getGitSource() string {
	gs := model.GitSource{}
	return model.GetGitSourceString(gs)
}

func getBuildConfiguration() string {
	bc := model.BuildConfiguration{}
	return model.GetBuildConfigurationString(bc)
}

func getCreator() string {
	c := model.Actor{}
	return model.GetCreatorString(c)
}

func withPrefix(body string, prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, body)
}

func getState() string {
	states := []string{
		string(model.DeploymentStateQueueing),
		string(model.DeploymentStateBuilding),
		string(model.DeploymentStateBuildSucceeded),
		string(model.DeploymentStateCommitted),
		string(model.DeploymentStateReady),
		string(model.DeploymentStateError),
	}
	return states[rand.Intn(6)]
}

func getType() string {
	states := []string{
		string(model.INFO),
		string(model.DEBUG),
		string(model.ERROR),
	}
	return states[rand.Intn(3)]
}
