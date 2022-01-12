package main

import "github.com/Jeffail/gabs/v2"

type Rate struct {
	Email string
	Year  int
	Level int8
}

func (a *Rate) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.Year, "year")
	_, _ = res.Set(a.Level, "level")
	return res
}
