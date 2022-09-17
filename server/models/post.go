package models

import (
	"github.com/Jeffail/gabs/v2"
	"strconv"
)

type Post struct {
	ID         uint32 `gorm:"autoIncrement;primaryKey"`
	Link       string `gorm:"not null"`
	Title      string `gorm:"not null"`
	Headline   string `gorm:"not null"`
	Content    string `gorm:"not null"`
	Hashtag    string
	Privacy    uint8
	LikeCount  uint
	ViewCount  uint
	UpdateDate uint64 `gorm:"autoUpdateTime:milli"`
	CreateDate uint64 `gorm:"autoCreateTime:milli"`
}

func (p *Post) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(p.ID, "id")
	_, _ = res.Set(p.Link+"."+strconv.FormatUint(uint64(p.ID), 10), "link")
	_, _ = res.Set(p.Title, "title")
	//_, _ = res.Set(p.Content, "content")
	_, _ = res.Set(p.Privacy, "privacy")
	_, _ = res.Set(p.Headline, "headline")
	_, _ = res.Set(p.Hashtag, "hashtag")
	_, _ = res.Set(p.ViewCount, "stats.views")
	_, _ = res.Set(p.LikeCount, "stats.likes")
	_, _ = res.Set(p.UpdateDate, "updateDate")
	_, _ = res.Set(p.CreateDate, "createDate")
	return res
}
