package localizer

import (
	"ffmpegGui/setting"
)

type RepositoryContract interface {
	GetCode() (string, error)
	Save(code string) (setting.Setting, error)
}

type Repository struct {
	settingRepository setting.RepositoryContract
}

func NewRepository(settingRepository setting.RepositoryContract) *Repository {
	return &Repository{settingRepository: settingRepository}
}

func (r Repository) GetCode() (string, error) {
	return r.settingRepository.GetValue("language")
}

func (r Repository) Save(code string) (setting.Setting, error) {
	return r.settingRepository.CreateOrUpdate("language", code)
}
