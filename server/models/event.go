package models

import (
	"github.com/Jeffail/gabs/v2"
	"strconv"
	"time"
)

type Event struct {
	ID         uint32 `gorm:"autoIncrement;primaryKey"`
	Link       string
	Title      string
	BeginDate  uint64
	EndDate    uint64
	Privacy    uint8
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (e *Event) GetStatus() string {
	now := uint64(time.Now().UnixMilli())
	if now >= e.BeginDate && now <= e.EndDate {
		return "ongoing"
	} else if now > e.EndDate {
		return "finished"
	} else {
		return "waiting"
	}
}

func (e *Event) Serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(e.ID, "id")
	_, _ = res.Set(e.Link+"."+strconv.FormatUint(uint64(e.ID), 10), "link")
	_, _ = res.Set(e.Title, "title")
	_, _ = res.Set(e.BeginDate, "beginDate")
	_, _ = res.Set(e.EndDate, "endDate")
	_, _ = res.Set(e.GetStatus(), "status")
	_, _ = res.Set(e.Privacy, "privacy")
	_, _ = res.Set(e.UpdateDate, "updateDate")
	_, _ = res.Set(e.CreateDate, "createDate")
	return res
}
