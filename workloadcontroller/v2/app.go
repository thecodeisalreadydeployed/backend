package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
)

func (ctrl *workloadController) NewApp(app *model.App) (*model.App, error) {
	a, createErr := datastore.SaveApp(datastore.GetDB(), app)
	if createErr != nil {
		return nil, createErr
	}

	err := ctrl.gitOpsController.SetupApp(a.ProjectID, a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}
