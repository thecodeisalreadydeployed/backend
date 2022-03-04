package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
)

func (ctrl *workloadController) NewProject(project *model.Project, dataStore datastore.DataStore) (*model.Project, error) {
	prj, createErr := dataStore.SaveProject(project)
	if createErr != nil {
		return nil, createErr
	}

	err := ctrl.gitOpsController.SetupProject(prj.ID)
	if err != nil {
		return nil, err
	}

	return prj, nil
}
