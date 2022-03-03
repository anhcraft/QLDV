package main

import "github.com/Jeffail/gabs/v2"

type Attachment struct {
	ID     string `gorm:"primaryKey"`
	Date   int64
	PostId string `gorm:"primaryKey"`
	Post   Post   `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a *Attachment) serialize() interface{} {
	res := gabs.New()
	_, _ = res.Set(a.ID, "id")
	//_, _ = res.Set(a.PostId, "postId")
	_, _ = res.Set(a.Date, "date")
	return res
}
