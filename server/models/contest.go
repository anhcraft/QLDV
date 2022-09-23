package models

import "github.com/Jeffail/gabs/v2"

type Contest struct {
	AcceptingAnswers bool
	LimitQuestions   uint8
	LimitTime        uint32
	LimitSessions    uint8
	DataSheet        string
	Info             string
	EventID          int   `gorm:"primaryKey"`
	Event            Event `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a *Contest) serialize(showData bool) *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.EventID, "id")
	_, _ = res.Set(a.AcceptingAnswers, "acceptingAnswers")
	_, _ = res.Set(a.LimitQuestions, "limitQuestions")
	_, _ = res.Set(a.LimitTime, "limitTime")
	_, _ = res.Set(a.LimitSessions, "limitSessions")
	if showData {
		_, _ = res.Set(a.DataSheet, "dataSheet")
	}
	_, _ = res.Set(a.Info, "info")
	return res
}
