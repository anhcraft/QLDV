package main

import (
	"github.com/Jeffail/gabs/v2"
	"strconv"
)

type Post struct {
	ID        int `gorm:"autoIncrement,primaryKey"`
	Link      string
	Title     string
	Headline  string
	Content   string
	Hashtag   string
	Date      int64
	Privacy   uint8
	LikeCount uint
	ViewCount uint
}

func (p *Post) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(p.ID, "id")
	_, _ = res.Set(p.Link+"."+strconv.Itoa(p.ID), "link")
	_, _ = res.Set(p.Title, "title")
	//_, _ = res.Set(p.Content, "content")
	_, _ = res.Set(p.Date, "date")
	_, _ = res.Set(p.Privacy, "privacy")
	_, _ = res.Set(p.Headline, "headline")
	_, _ = res.Set(p.Hashtag, "hashtag")
	_, _ = res.Set(p.ViewCount, "views")
	_, _ = res.Set(p.LikeCount, "likes")
	return res
}
