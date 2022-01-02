package main

import "github.com/Jeffail/gabs/v2"

type Achievement struct {
	Email string
	Title string
	Year int
}

func (a *Achievement) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.Title, "title")
	_, _ = res.Set(a.Year, "year")
	return res
}
