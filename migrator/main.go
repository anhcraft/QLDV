package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
	"time"
)

func main() {
	f, err := excelize.OpenFile("dv.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows("KHOI 12")
	if err != nil {
		fmt.Println(err)
		return
	}
	users := make([]*User, 0)
	for i, row := range rows {
		if i < 3 {
			continue
		}
		if i > 537 {
			break
		}
		//fmt.Println("Processing ", i, " (", row[2], row[3], ")...")
		birth, err := time.Parse("02/01/2006", row[5])
		if err != nil {
			fmt.Println(err)
			return
		}
		gender := strings.Contains(row[4], "Nữ")
		if !gender && !strings.Contains(row[4], "Nam") {
			fmt.Println("Invalid gender: ", row[4])
			return
		}
		user := &User{
			Class:   "12" + strings.TrimSpace(row[1]),
			Name:    ToCamel(strings.TrimSpace(row[2]) + " " + strings.TrimSpace(row[3])),
			Gender:  gender,
			Birth:   birth.UnixMilli(),
			Address: strings.TrimSpace(row[12]),
			Email:   strings.TrimSpace(row[8]),
		}
		if len(row) == 15 {
			fmt.Println("DV Status not found: ")
			user.Certified = false
		} else {
			if strings.ToLower(strings.TrimSpace(row[15])) == "đoàn viên" || strings.ToLower(strings.TrimSpace(row[15])) == "đoan viên" || strings.ToLower(strings.TrimSpace(row[15])) == "đoàn viên" || strings.ToLower(strings.TrimSpace(row[15])) == "doàn viên" || strings.ToUpper(strings.TrimSpace(row[15])) == "DV" || strings.ToUpper(strings.TrimSpace(row[15])) == "ĐV" {
				user.Certified = true
			} else {
				fmt.Println("DV Status not found: ", row[15])
				user.Certified = false
			}
		}
		users = append(users, user)
	}
	js, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile("./12.json", js, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ToCamel(s string) string {
	s = strings.ToLower(s)
	var g []string
	p := strings.Fields(s)
	for _, value := range p {
		g = append(g, strings.Title(value))
	}
	return strings.Join(g, " ")
}
