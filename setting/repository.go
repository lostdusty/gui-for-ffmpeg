package setting

import (
	"errors"
	"gorm.io/gorm"
)

type RepositoryContract interface {
	Create(setting Setting) (Setting, error)
	CreateOrUpdate(code string, value string) (Setting, error)
	GetValue(code string) (value string, err error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) GetValue(code string) (value string, err error) {
	var setting Setting
	err = r.db.Where("code = ?", code).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, err
}

func (r Repository) Create(setting Setting) (Setting, error) {
	err := r.db.Create(&setting).Error
	if err != nil {
		return setting, err
	}
	return setting, err
}

func (r Repository) CreateOrUpdate(code string, value string) (Setting, error) {
	var setting Setting
	err := r.db.Where("code = ?", code).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == true {
			setting = Setting{Code: code, Value: value}
			return r.Create(setting)
		} else {
			return setting, err
		}
	}
	err = r.db.Model(&setting).UpdateColumn("value", value).Error
	if err != nil {
		return setting, err
	}
	return setting, err
}
