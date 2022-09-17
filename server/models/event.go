package models

import (
	"github.com/Jeffail/gabs/v2"
	"strconv"
	"time"
)

type Event struct {
	ID        int `gorm:"autoIncrement;primaryKey"`
	Title     string
	Link      string
	BeginDate int64
	EndDate   int64
	Date      int64
	Privacy   uint8
}

func (e *Event) GetStatus() string {
	now := time.Now().UnixMilli()
	if now >= e.BeginDate && now <= e.EndDate {
		return "ongoing"
	} else if now > e.EndDate {
		return "finished"
	} else {
		return "waiting"
	}
}

func (e *Event) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(e.ID, "id")
	_, _ = res.Set(e.Title, "title")
	_, _ = res.Set(e.Link+"."+strconv.Itoa(e.ID), "link")
	_, _ = res.Set(e.BeginDate, "beginDate")
	_, _ = res.Set(e.EndDate, "endDate")
	_, _ = res.Set(e.Date, "date")
	_, _ = res.Set(e.GetStatus(), "status")
	_, _ = res.Set(e.Privacy, "privacy")
	return res
}
