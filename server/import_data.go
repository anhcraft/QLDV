package main

import (
	"das/models"
	"das/storage"
	"das/utils"
	"encoding/json"
	"io/ioutil"
)

type UserD struct {
	Class     string `json:"class"`
	Name      string `json:"name"`
	Birth     int64  `json:"birth"`
	Gender    bool   `json:"gender"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Certified bool   `json:"certified"`
}

func importData() {
	usersD := make([]*UserD, 0)

	content, _ := ioutil.ReadFile("12.json")

	_ = json.Unmarshal(content, &usersD)

	users := make([]models.User, 0)
	for _, v := range usersD {
		r := utils.RoleRegularMember
		if v.Certified {
			r = utils.RoleCertifiedMember
		}
		users = append(users, models.User{
			Email:     v.Email,
			Name:      v.Name,
			Gender:    v.Gender,
			Birthday:  uint64(v.Birth),
			EntryYear: 2019,
			Phone:     "",
			Role:      r,
			Class:     v.Class,
		})
	}

	storage.Db.CreateInBatches(users, 20)
}
