package main

import "github.com/Jeffail/gabs/v2"

type Rate struct {
	Email string
	Year int
	Excellent bool
}

func (a *Rate) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.Year, "year")
	_, _ = res.Set(a.Excellent, "good")
	return res
}
