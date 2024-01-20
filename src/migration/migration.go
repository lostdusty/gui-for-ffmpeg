package migration

import (
	"ffmpegGui/setting"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	err := db.AutoMigrate(&setting.Setting{})
	if err != nil {
		return err
	}

	return nil
}
