package models

type PostStat struct {
	PostId     uint32 `gorm:"primaryKey"`
	Post       Post   `gorm:"constraint:OnDelete:CASCADE;"`
	UserId     uint16 `gorm:"primaryKey"`
	Action     string `gorm:"primaryKey"`
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (PostStat) TableName() string {
	return "post_stats"
}
