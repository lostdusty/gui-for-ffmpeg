package convertor

import (
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/src/setting"
)

type RepositoryContract interface {
	GetPathFfmpeg() (string, error)
	SavePathFfmpeg(code string) (setting.Setting, error)
	GetPathFfprobe() (string, error)
	SavePathFfprobe(code string) (setting.Setting, error)
}

type Repository struct {
	settingRepository setting.RepositoryContract
}

func NewRepository(settingRepository setting.RepositoryContract) *Repository {
	return &Repository{settingRepository: settingRepository}
}

func (r Repository) GetPathFfmpeg() (string, error) {
	return r.settingRepository.GetValue("ffmpeg")
}

func (r Repository) SavePathFfmpeg(path string) (setting.Setting, error) {
	return r.settingRepository.CreateOrUpdate("ffmpeg", path)
}

func (r Repository) GetPathFfprobe() (string, error) {
	return r.settingRepository.GetValue("ffprobe")
}

func (r Repository) SavePathFfprobe(path string) (setting.Setting, error) {
	return r.settingRepository.CreateOrUpdate("ffprobe", path)
}
