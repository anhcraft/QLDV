package main

import "github.com/Jeffail/gabs/v2"

type PostStat struct {
	PostId string `gorm:"primaryKey"`
	UserId string `gorm:"primaryKey"`
	Action string `gorm:"primaryKey"`
	Date   int64
}

func (PostStat) TableName() string {
	return "post_stats"
}

func (p *PostStat) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(p.PostId, "postId")
	_, _ = res.Set(p.UserId, "userId")
	_, _ = res.Set(p.Action, "action")
	_, _ = res.Set(p.Date, "date")
	return res
}
