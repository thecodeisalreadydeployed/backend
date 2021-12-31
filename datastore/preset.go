package datastore

import (
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/gorm"
)

func GetAllPresets(DB *gorm.DB) (*[]model.Preset, error) {
	return nil, nil
}

func GetPresetByID(DB *gorm.DB, id string) (*model.Preset, error) {
	return nil, nil
}

func GetPresetByName(DB *gorm.DB, name string) (*model.Preset, error) {
	return nil, nil
}

func SavePreset(DB *gorm.DB, preset *model.Preset) (*model.Preset, error) {
	return nil, nil
}

func RemovePreset(DB *gorm.DB, id string) error {
	return nil
}
