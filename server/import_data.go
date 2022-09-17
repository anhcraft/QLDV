package main

import (
	"das/models"
	"encoding/json"
	"fmt"
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
	setupDB()

	usersD := make([]*UserD, 0)

	content, _ := ioutil.ReadFile("12.json")

	_ = json.Unmarshal(content, &usersD)

	users := make([]models.User, 0)
	for i, v := range usersD {
		users = append(users, models.User{
			Email:     v.Email,
			StudentId: fmt.Sprintf("%016d", 9099_0012_1922_0000+i+1),
			Name:      v.Name,
			Gender:    v.Gender,
			Birthday:  v.Birth,
			EntryYear: 2019,
			Phone:     "",
			Certified: v.Certified,
			Class:     v.Class,
			Admin:     false,
		})
	}

	Db.CreateInBatches(users, 20)
}
