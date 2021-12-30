package datastore

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func getDeploymentRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.DeploymentStructString()).
		AddRow(
			"dpl-test",
			"app-test",
			model.GetCreatorString(model.Actor{}),
			"dummy_meta",
			model.GetGitSourceString(model.GitSource{}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.DeploymentStateReady,
		)
}

func getAppRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.AppStructString()).
		AddRow(
			"app-test",
			"prj-test",
			"Best App",
			model.GetGitSourceString(model.GitSource{}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			true,
		)
}

func getProjectRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.ProjectStructString()).
		AddRow("prj-test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
}

func getExpectedDeployment() *model.Deployment {
	return &model.Deployment{
		ID:                 "dpl-test",
		AppID:              "app-test",
		Creator:            model.Actor{},
		Meta:               "dummy_meta",
		GitSource:          model.GitSource{},
		BuiltAt:            time.Unix(0, 0),
		CommittedAt:        time.Unix(0, 0),
		DeployedAt:         time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		State:              model.DeploymentStateReady,
	}
}

func getExpectedApp() *model.App {
	return &model.App{
		ID:                 "app-test",
		ProjectID:          "prj-test",
		Name:               "Best App",
		GitSource:          model.GitSource{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         true,
	}
}

func getExpectedProject() *model.Project {
	return &model.Project{
		ID:        "prj-test",
		Name:      "Best Project",
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}
}
