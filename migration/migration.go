package migration

import (
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/setting"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	err := db.AutoMigrate(&setting.Setting{})
	if err != nil {
		return err
	}

	return nil
}
