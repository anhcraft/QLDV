package models

import (
	"github.com/Jeffail/gabs/v2"
)

const UnknownRank = 0
const ExcellentRank = 1
const GoodRank = 2
const MediumRank = 3

type AnnualRank struct {
	UserId     uint16 `gorm:"primaryKey"`
	User       User   `gorm:"constraint:OnDelete:CASCADE"`
	Year       uint16
	Level      uint8
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (a *AnnualRank) Serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.Year, "year")
	_, _ = res.Set(a.Level, "level")
	_, _ = res.Set(a.UpdateDate, "updateDate")
	_, _ = res.Set(a.CreateDate, "createDate")
	return res
}
