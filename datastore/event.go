package datastore

import (
	"errors"
	"regexp"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/gorm"
)

func (d *dataStore) GetEventsByDeploymentID(deploymentID string) (*[]model.Event, error) {
	var _data []datamodel.Event
	err := d.DB.Table("events").Where(datamodel.Event{DeploymentID: deploymentID}).Scan(&_data).Error

	if err != nil {
		d.logger.Error(err.Error())
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

func (d *dataStore) GetEventByID(eventID string) (*model.Event, error) {
	if !(d.IsValidKSUID(eventID)) {
		d.logger.Error(MsgEventPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.Event
	err := d.DB.First(&_data, "id = ?", eventID).Error

	if err != nil {
		d.logger.Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func (d *dataStore) SaveEvent(event *model.Event) (*model.Event, error) {
	if event.ID == "" {
		event.ID = model.GenerateEventID(event.ExportedAt)
	}
	if !(d.IsValidKSUID(event.ID)) {
		d.logger.Error(MsgEventPrefix)
		return nil, errutil.ErrInvalidArgument
	}
	e := datamodel.NewEventFromModel(event)
	err := d.DB.Save(e).Error

	if err != nil {
		d.logger.Error(err.Error())

		if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrInvalidData) {
			return nil, errutil.ErrInvalidArgument
		}

		if errors.Is(err, gorm.ErrMissingWhereClause) {
			return nil, errutil.ErrFailedPrecondition
		}

		return nil, errutil.ErrUnknown
	}

	return d.GetEventByID(event.ID)
}

func (d *dataStore) IsValidKSUID(str string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]{27}$")
	return re.MatchString(str)
}
