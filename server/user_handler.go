package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

func getProfile(email string) *User {
	var user User
	result := db.Take(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &user
	}
}

func getAchievements(email string) []Achievement {
	var achievements []Achievement
	_ = db.Find(&achievements, "email = ?", email)
	return achievements
}

func getRates(email string) []Rate {
	var rates []Rate
	_ = db.Find(&rates, "email = ?", email)
	return rates
}

func saveAchievement(achievements []struct {
	Title string `json:"title,omitempty"`
	Year  int    `json:"year,omitempty"`
}, user string) {
	var ach Achievement
	db.Where("email = ?", user).Delete(&ach)

	_achievements := make([]Achievement, 0)
	for _, v := range achievements {
		if len(strings.TrimSpace(v.Title)) == 0 {
			continue
		}
		_achievements = append(_achievements, Achievement{
			Email: user,
			Title: v.Title,
			Year:  v.Year,
		})
	}
	if len(_achievements) > 0 {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(_achievements)
	}
}

func saveRates(rates map[int]int8, user string) {
	var rate Rate
	db.Where("email = ?", user).Delete(&rate)

	_rates := make([]Rate, 0)
	for k, v := range rates {
		_rates = append(_rates, Rate{
			Email: user,
			Year:  k,
			Level: v,
		})
	}
	if len(_rates) > 0 {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(_rates)
	}
}

func getUsers(limit int, offset int, name string, class string, email string, certified int) []User {
	name = strings.ToLower(strings.TrimSpace(name))
	class = strings.ToLower(strings.TrimSpace(class))
	email = strings.ToLower(strings.TrimSpace(email))

	var users []User
	a := db.Offset(offset).Order("student_id").Limit(limit)
	if len(name) > 0 {
		a = a.Where("LOWER(`name`) like ?", "%"+name+"%")
	}
	if len(class) > 0 {
		a = a.Where("LOWER(`class`) like ?", "%"+class+"%")
	}
	if len(email) > 0 {
		a = a.Where("LOWER(`email`) like ?", "%"+email+"%")
	}
	if certified == 1 {
		a = a.Where("`certified` = 1")
	} else if certified == 2 {
		a = a.Where("`certified` = 0")
	}
	a = a.Find(&users)
	return users
}

// TODO merge user changes to solve n+1 problem

func setUserCertified(certified map[string]bool, class string) bool {
	var user User
	for k, v := range certified {
		a := db.Model(&user).Where("email = ?", k)
		if class != "" {
			a = a.Where("class = ?", class)
		}
		a.Update("certified", v)
	}
	return true
}

func setUserModStatus(mod map[string]bool) bool {
	var user User
	for k, v := range mod {
		db.Model(&user).Where("email = ?", k).Update("mod", v)
	}
	return true
}

func profileGetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	user := getProfile(email)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	return c.SendString(user.serialize().String())
}

func progressionGetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}

	payload := struct {
		User string `json:"user,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	if payload.User == "" {
		payload.User = email
	} else if payload.User != email {
		user := getProfile(email)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		if !user.Admin && !user.Mod {
			_, _ = res.Set("ERR_NO_PERMISSION", "error")
			return c.SendString(res.String())
		}
	}

	_, _ = res.Array("rates")
	for _, rate := range getRates(payload.User) {
		_ = res.ArrayAppend(rate.serialize(), "rates")
	}
	_, _ = res.Array("achievements")
	for _, rate := range getAchievements(payload.User) {
		_ = res.ArrayAppend(rate.serialize(), "achievements")
	}
	return c.SendString(res.String())
}

func progressionChangeRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	requester := getProfile(email)
	if requester == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !requester.Admin && !requester.Mod {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	payload := struct {
		User         string       `json:"user,omitempty"`
		Rates        map[int]int8 `json:"rates,omitempty"`
		Achievements []struct {
			Title string `json:"title,omitempty"`
			Year  int    `json:"year,omitempty"`
		} `json:"achievements,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	saveRates(payload.Rates, payload.User)
	saveAchievement(payload.Achievements, payload.User)
	_, _ = res.Set(true, "success")
	return c.SendString(res.String())
}

func userListRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	requester := getProfile(email)
	if requester == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !requester.Admin && !requester.Mod {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	payload := struct {
		Limit           int    `json:"limit,omitempty"`
		Offset          int    `json:"offset,omitempty"`
		FilterName      string `json:"filter_name,omitempty"`
		FilterClass     string `json:"filter_class,omitempty"`
		FilterEmail     string `json:"filter_email,omitempty"`
		FilterCertified int    `json:"filter_certified,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	if requester.Mod {
		payload.FilterClass = requester.Class
	}

	if payload.Limit > 50 {
		payload.Limit = 50
	}
	if payload.Offset < 0 {
		payload.Offset = 0
	}
	_, _ = res.Array("users")
	for _, post := range getUsers(payload.Limit, payload.Offset, payload.FilterName, payload.FilterClass, payload.FilterEmail, payload.FilterCertified) {
		_ = res.ArrayAppend(post.serialize(), "users")
	}
	return c.SendString(res.String())
}

func userChangeRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	requester := getProfile(email)
	if requester == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !requester.Admin && !requester.Mod {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	payload := struct {
		Certified map[string]bool `json:"certified,omitempty"`
		Mod       map[string]bool `json:"mod,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	class := ""
	if requester.Mod {
		class = requester.Class
	}

	_, _ = res.Set(setUserCertified(payload.Certified, class), "success")

	if requester.Admin {
		_, _ = res.Set(setUserModStatus(payload.Mod), "success")
	}

	return c.SendString(res.String())
}

func userStatGetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	requester := getProfile(email)
	if requester == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !requester.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	result := struct {
		a int64
		b int64
		c int64
		d int64
		e int64
	}{}

	x := db.Raw("select count(if(gender = true and (class like '10%' or class like '11%' or class like '12%'), 1, null)) as a, count(if(certified = true and (class like '10%' or class like '11%' or class like '12%'), 1, null)) as b, count(if(class like '10%', 1, null)) as c, count(if(class like '11%', 1, null)) as d, count(if(class like '12%', 1, null)) as e from users")
	_ = x.Row().Scan(&result.a, &result.b, &result.c, &result.d, &result.e)
	_, _ = res.Set(result.a, "women")
	_, _ = res.Set(result.b, "certified")
	_, _ = res.Set(result.c, "class10")
	_, _ = res.Set(result.d, "class11")
	_, _ = res.Set(result.e, "class12")
	return c.SendString(res.String())
}

func profileCoverSetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	user := getProfile(email)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}

	t := c.Get("content-type")

	if t == "image/png" {
		_, _ = res.Set(setProfileCover(email, c.Body(), ".png"), "success")
	} else if t == "image/jpeg" {
		_, _ = res.Set(setProfileCover(email, c.Body(), ".jpeg"), "success")
	} else {
		_, _ = res.Set("ERR_ILLEGAL_PROFILE_COVER", "error")
		return c.SendString(res.String())
	}

	return c.SendString(res.String())
}

func setProfileCover(email string, data []byte, ext string) bool {
	_ = os.Mkdir("public", os.ModePerm)
	hash := sha256.New()
	hash.Write([]byte(email))
	md := hash.Sum(nil)
	id := "cover-" + hex.EncodeToString(md)
	err := os.WriteFile("./public/"+id+ext, data, os.ModePerm)
	if err != nil {
		return false
	}
	db.Model(&User{}).Where("email = ?", email).Update("profile_cover", id+ext)
	return true
}
