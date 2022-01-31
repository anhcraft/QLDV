package main

import "github.com/Jeffail/gabs/v2"

type Event struct {
	ID        string
	Title     string
	StartDate int64
	EndDate   int64
	Date      int64
}

func (e *Event) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(e.ID, "id")
	_, _ = res.Set(e.Title, "title")
	_, _ = res.Set(e.StartDate, "startDate")
	_, _ = res.Set(e.EndDate, "endDate")
	_, _ = res.Set(e.Date, "date")
	return res
}
