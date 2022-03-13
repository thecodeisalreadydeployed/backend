package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
)

func (ctrl *workloadController) NewApp(app *model.App, dataStore datastore.DataStore) (*model.App, error) {
	a, createErr := dataStore.SaveApp(app)
	if createErr != nil {
		return nil, createErr
	}

	err := ctrl.gitOpsController.SetupApp(a.ProjectID, a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}
