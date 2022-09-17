package models

import (
	"github.com/Jeffail/gabs/v2"
)

type Attachment struct {
	ID         string `gorm:"primaryKey"`
	PostId     uint32 `gorm:"primaryKey"`
	Post       Post   `gorm:"constraint:OnDelete:CASCADE;"`
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (a *Attachment) serialize() interface{} {
	res := gabs.New()
	_, _ = res.Set(a.ID, "id")
	_, _ = res.Set(a.UpdateDate, "updateDate")
	_, _ = res.Set(a.CreateDate, "createDate")
	return res
}
