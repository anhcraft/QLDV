package main

import "github.com/Jeffail/gabs/v2"

type Contest struct {
	AcceptingAnswers bool
	LimitQuestions   uint16
	LimitTime        uint32
	DataSheet        string
	Info             string
	EventID          string `gorm:"primaryKey"`
	Event            Event
}

func (a *Contest) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.EventID, "id")
	_, _ = res.Set(a.AcceptingAnswers, "acceptingAnswers")
	_, _ = res.Set(a.LimitQuestions, "limitQuestions")
	_, _ = res.Set(a.LimitTime, "limitTime")
	_, _ = res.Set(a.DataSheet, "dataSheet")
	_, _ = res.Set(a.Info, "info")
	return res
}
