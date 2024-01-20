package setting

type Setting struct {
	ID    uint   `gorm:"primary_key"`
	Code  string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Value string `gorm:"type:text"`
}
