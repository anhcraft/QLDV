package utils

import (
	"github.com/rs/zerolog/log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func RemoveVietnameseAccents(str string) string {
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
	str = reg_a.ReplaceAllLiteralString(str, "a")
	str = reg_e.ReplaceAllLiteralString(str, "e")
	str = reg_i.ReplaceAllLiteralString(str, "i")
	str = reg_o.ReplaceAllLiteralString(str, "o")
	str = reg_u.ReplaceAllLiteralString(str, "u")
	str = reg_y.ReplaceAllLiteralString(str, "y")
	str = reg_d.ReplaceAllLiteralString(str, "d")
	var RegexpPara = `\(.*\)`
	regPara := regexp.MustCompile(RegexpPara)
	str = regPara.ReplaceAllLiteralString(str, "")
	return strings.ToLower(str)
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

func GenerateLinkFromTitle(title string) string {
	id := strings.ToLower(strings.TrimSpace(title))
	id = RemoveVietnameseAccents(id)
	id = regexp.MustCompile("[^a-z0-9 ]+").ReplaceAllString(id, "")
	id = strings.ReplaceAll(id, " ", "-")
	return id
}

// ValidateNonNegativeInteger Checks a string is an integer where N >= 0
func ValidateNonNegativeInteger(str string) bool {
	b, err := regexp.MatchString("^[0-9]+$", str)
	if err != nil {
		log.Error().Err(err).Str("str", str).Msg("An error occurred while validating integer")
		return false
	}
	return b
}

func ClampUint8(v uint8, min uint8, max uint8) uint8 {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}
