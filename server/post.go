package main

import "github.com/Jeffail/gabs/v2"

type Post struct {
	ID      string
	Title   string
	Content string
	Date    int64
	Privacy uint8
}

func (p *Post) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(p.ID, "id")
	_, _ = res.Set(p.Title, "title")
	//_, _ = res.Set(p.Content, "content")
	_, _ = res.Set(p.Date, "date")
	_, _ = res.Set(p.Privacy, "privacy")
	return res
}
