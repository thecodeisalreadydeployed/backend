package datastore

import (
	"regexp"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetEventsByDeploymentID(DB *gorm.DB, deploymentID string) (*[]model.Event, error) {
	var _data []datamodel.Event
	err := DB.Table("events").Where(datamodel.Event{DeploymentID: deploymentID}).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.Event
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
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
		event.ID = model.GenerateEventID(event.ExportedAt)
	}
	e := datamodel.NewEventFromModel(event)
	err := DB.Save(e).Error

	if err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrUnknown
	}

	return nil
}
