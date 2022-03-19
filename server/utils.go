package main

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func RemoveVietnameseAccents(province string) string {
	var RegexpA = `à|á|ạ|ã|ả|ă|ắ|ằ|ẳ|ẵ|ặ|â|ấ|ầ|ẩ|ẫ|ậ`
	var RegexpE = `è|ẻ|ẽ|é|ẹ|ê|ề|ể|ễ|ế|ệ`
	var RegexpI = `ì|ỉ|ĩ|í|ị`
	var RegexpU = `ù|ủ|ũ|ú|ụ|ư|ừ|ử|ữ|ứ|ự`
	var RegexpY = `ỳ|ỷ|ỹ|ý|ỵ`
	var RegexpO = `ò|ỏ|õ|ó|ọ|ô|ồ|ổ|ỗ|ố|ộ|ơ|ờ|ở|ỡ|ớ|ợ`
	var RegexpD = `Đ|đ`
	reg_a := regexp.MustCompile(RegexpA)
	reg_e := regexp.MustCompile(RegexpE)
	reg_i := regexp.MustCompile(RegexpI)
	reg_o := regexp.MustCompile(RegexpO)
	reg_u := regexp.MustCompile(RegexpU)
	reg_y := regexp.MustCompile(RegexpY)
	reg_d := regexp.MustCompile(RegexpD)
	province = reg_a.ReplaceAllLiteralString(province, "a")
	province = reg_e.ReplaceAllLiteralString(province, "e")
	province = reg_i.ReplaceAllLiteralString(province, "i")
	province = reg_o.ReplaceAllLiteralString(province, "o")
	province = reg_u.ReplaceAllLiteralString(province, "u")
	province = reg_y.ReplaceAllLiteralString(province, "y")
	province = reg_d.ReplaceAllLiteralString(province, "d")
	var RegexpPara = `\(.*\)`
	regPara := regexp.MustCompile(RegexpPara)
	province = regPara.ReplaceAllLiteralString(province, "")
	return strings.ToLower(province)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateIdFromTitle(title string) string {
	id := strings.ToLower(strings.TrimSpace(title))
	id = RemoveVietnameseAccents(id)
	id = regexp.MustCompile("[^a-z0-9 ]+").ReplaceAllString(id, "")
	id = strings.ReplaceAll(id, " ", "-")
	id += "-" + strings.ToLower(RandStringBytes(5))
	return id
}
