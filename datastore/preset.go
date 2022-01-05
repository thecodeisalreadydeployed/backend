package datastore

import (
	"errors"
	"fmt"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

func GetAllPresets(DB *gorm.DB) (*[]model.Preset, error) {
	var _data []datamodel.Preset
	err := DB.Table("presets").Scan(&_data).Error
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.Preset
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func GetPresetByID(DB *gorm.DB, presetID string) (*model.Preset, error) {
	if !strings.HasPrefix(presetID, "pst-") {
		zap.L().Error(MsgPresetPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.Preset
	err := DB.First(&_data, "id = ?", presetID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func GetPresetsByName(DB *gorm.DB, name string) (*[]model.Preset, error) {
	var _data []datamodel.Preset

	err := DB.Table("presets").Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.Preset
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func SavePreset(DB *gorm.DB, preset *model.Preset) (*model.Preset, error) {
	if preset.ID == "" {
		preset.ID = model.GeneratePresetID()
	}
	if !strings.HasPrefix(preset.ID, "pst-") {
		zap.L().Error(MsgPresetPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	a := datamodel.NewPresetFromModel(preset)
	err := DB.Save(a).Error

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
	return GetPresetByID(DB, preset.ID)
}

func RemovePreset(DB *gorm.DB, id string) error {
	if !strings.HasPrefix(id, "pst-") {
		zap.L().Error(MsgPresetPrefix)
		return errutil.ErrInvalidArgument
	}
	var a datamodel.Preset
	err := DB.Table("presets").Where(datamodel.Preset{ID: id}).First(&a).Error
	if err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrNotFound
	}
	if err := DB.Delete(&a).Error; err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrUnknown
	}
	return nil
}
