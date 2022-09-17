package models

import "github.com/Jeffail/gabs/v2"

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

func (p *PostStat) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(p.PostId, "postId")
	_, _ = res.Set(p.UserId, "userId")
	_, _ = res.Set(p.Action, "action")
	_, _ = res.Set(p.UpdateDate, "updateDate")
	_, _ = res.Set(p.CreateDate, "createDate")
	return res
}
