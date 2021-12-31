package datastore

import (
	"errors"
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

func GetEventByID(DB *gorm.DB, eventID string) (*model.Event, error) {
	if !IsValidKSUID(eventID) {
		zap.L().Error(MsgEventPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.Event
	err := DB.First(&_data, "id = ?", eventID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SaveEvent(DB *gorm.DB, event *model.Event) (*model.Event, error) {
	if event.ID == "" {
		event.ID = model.GenerateEventID(event.ExportedAt)
	}
	if !IsValidKSUID(event.ID) {
		zap.L().Error(MsgEventPrefix)
		return nil, errutil.ErrInvalidArgument
	}
	e := datamodel.NewEventFromModel(event)
	err := DB.Save(e).Error

	if err != nil {
		zap.L().Error(err.Error())

		if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrInvalidData) {
			return nil, errutil.ErrInvalidArgument
		}

		if errors.Is(err, gorm.ErrMissingWhereClause) {
			return nil, errutil.ErrFailedPrecondition
		}

		return nil, errutil.ErrUnknown
	}

	return GetEventByID(DB, event.ID)
}

func IsValidKSUID(str string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]{27}$")
	return re.MatchString(str)
}
