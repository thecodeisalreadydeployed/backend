package datastore

import (
	"regexp"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetEventByDeploymentID(DB *gorm.DB, deploymentID string) (string, error) {
	return "", nil
}

func IsValidKSUID(str string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]{27}$")
	return re.MatchString(str)
}

func SaveEvent(DB *gorm.DB, event *model.Event) error {
	if event.ID != "" {
		if !IsValidKSUID(event.ID) {
			return errutil.ErrInvalidArgument
		}
	} else {
		event.ID = model.GenerateEventID()
	}
	e := datamodel.NewEventFromModel(event)
	err := DB.Save(e).Error

	if err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrUnknown
	}

	return nil
}
