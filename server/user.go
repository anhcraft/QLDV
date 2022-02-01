package main

import "github.com/Jeffail/gabs/v2"

type User struct {
	Email        string `gorm:"primaryKey"`
	StudentId    string
	Name         string
	Gender       bool
	Birthday     int64
	EntryYear    int
	Phone        string
	Certified    bool
	Class        string
	Admin        bool
	Mod          bool
	ProfileCover string
}

func (u *User) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(u.Email, "email")
	_, _ = res.Set(u.Name, "name")
	_, _ = res.Set(u.Gender, "gender")
	_, _ = res.Set(u.Birthday, "birth")
	_, _ = res.Set(u.EntryYear, "entry")
	_, _ = res.Set(u.Phone, "phone")
	_, _ = res.Set(u.Certified, "certified")
	_, _ = res.Set(u.Class, "class")
	_, _ = res.Set(u.StudentId, "sid")
	_, _ = res.Set(u.Admin, "admin")
	_, _ = res.Set(u.Mod, "mod")
	_, _ = res.Set(u.ProfileCover, "profileCover")
	return res
}
