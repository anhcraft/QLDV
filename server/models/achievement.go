package models

import (
	"github.com/Jeffail/gabs/v2"
)

type Achievement struct {
	UserId uint16 `gorm:"primaryKey"`
	User   User   `gorm:"constraint:OnDelete:CASCADE"`
	Title  string
	Year   uint16
	// Date stuff:
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (a *Achievement) Serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.Title, "title")
	_, _ = res.Set(a.Year, "year")
	_, _ = res.Set(a.UpdateDate, "updateDate")
	_, _ = res.Set(a.CreateDate, "createDate")
	return res
}
