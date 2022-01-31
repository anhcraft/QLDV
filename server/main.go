package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var client *auth.Client
var db *gorm.DB

func setupFirebase() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		_ = fmt.Errorf("error initializing app: %v", err)
		return
	}
	client_, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	client = client_
}

func setupDB() {
	dsn := os.Getenv("sql")
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting database: %v\n", err)
	}
	err = db_.AutoMigrate(&User{}, &Rate{}, &Achievement{}, &Post{}, &Attachment{}, &Event{})
	if err != nil {
		log.Fatalf("error migrating: %v\n", err)
	}
	db = db_
}

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

func getPost(id string) *Post {
	var post Post
	result := db.Take(&post, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &post
	}
}

func postChange(id string, title string, content string) *Post {
	if id == "" {
		hash := sha256.New()
		hash.Write([]byte(id + title + time.Now().String()))
		md := hash.Sum(nil)
		id = hex.EncodeToString(md)
	}
	post := Post{
		ID:      id,
		Title:   title,
		Content: content,
		Date:    time.Now().UnixMilli(),
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "content"}),
	}).Create(&post)
	return &post
}

func removePost(id string) bool {
	var post Post
	db.Where("id = ?", id).Delete(&post)
	return true
}

func uploadAttachment(postId string, data []byte, ext string) bool {
	_ = os.Mkdir("public", os.ModePerm)
	hash := sha256.New()
	hash.Write([]byte(postId + time.Now().String()))
	md := hash.Sum(nil)
	attId := hex.EncodeToString(md)
	att := Attachment{
		ID:     attId + ext,
		PostId: postId,
		Date:   time.Now().UnixMilli(),
	}
	err := os.WriteFile("./public/"+attId+ext, data, os.ModePerm)
	if err != nil {
		return false
	}
	_ = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&att)
	return true
}

func getAttachments(postId string) []Attachment {
	var atts []Attachment
	_ = db.Where("post_id = ?", postId).Order("date desc").Find(&atts)
	return atts
}

func removeAttachments(postId string, atts []string) bool {
	var att Attachment
	for _, v := range atts {
		db.Where("id = ?", v).Where("post_id = ?", postId).Delete(&att)
	}
	return true
}

func getPosts(limit int, older int64) []Post {
	var posts []Post
	_ = db.Where("date < ?", older).Order("date desc").Limit(limit).Find(&posts)
	return posts
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

func changeUserCert(certified map[string]bool, class string) bool {
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

func changeUserMod(mod map[string]bool) bool {
	var user User
	for k, v := range mod {
		db.Model(&user).Where("email = ?", k).Update("mod", v)
	}
	return true
}

func analyzeTokenToEmail(token string, c context.Context) (bool, string) {
	tkn, err := client.VerifyIDToken(c, token)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return false, "ERR_TOKEN_VERIFY"
	}
	u, err := client.GetUser(c, tkn.UID)
	if err != nil {
		log.Printf("error getting user %s: %v\n", tkn.UID, err)
		return false, "ERR_USER_GET"
	}
	return true, u.Email
}

func getEvents(limit int, older int64, fromDate int64, toDate int64) []Event {
	var events []Event
	a := db.Where("date < ?", older)
	if fromDate != 0 {
		a = a.Where("start_date >= ? or end_date >= ?", fromDate, fromDate)
	}
	if toDate != 0 {
		a = a.Where("start_date <= ? or end_date <= ?", toDate, toDate)
	}
	a.Order("date desc").Limit(limit).Find(&events)
	return events
}

func removeEvent(id string) bool {
	var event Event
	db.Where("id = ?", id).Delete(&event)
	return true
}

func getEvent(id string) *Event {
	var event Event
	result := db.Take(&event, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &event
	}
}

func editOrCreateEvent(id string, title string, startDate int64, endDate int64) *Event {
	if id == "" {
		hash := sha256.New()
		hash.Write([]byte(id + title + time.Now().String()))
		md := hash.Sum(nil)
		id = hex.EncodeToString(md)
	}
	event := Event{
		ID:        id,
		Title:     title,
		StartDate: startDate,
		EndDate:   endDate,
		Date:      time.Now().UnixMilli(),
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "start_date", "end_date"}),
	}).Create(&event)
	return &event
}

func main() {
	setupFirebase()
	setupDB()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Post("/profile", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		user := getProfile(txt)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		return c.SendString(user.serialize().String())
	})

	app.Post("/progression", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
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
			payload.User = txt
		} else if payload.User != txt {
			user := getProfile(txt)
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
	})

	app.Post("/change-progression", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		requester := getProfile(txt)
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
	})

	app.Get("/posts", func(c *fiber.Ctx) error {
		res := gabs.New()
		limit, err1 := strconv.Atoi(c.Query("limit", ""))
		if err1 != nil || limit > 50 {
			limit = 50
		}
		older, err2 := strconv.ParseInt(c.Query("older", ""), 10, 64)
		if err2 != nil {
			older = time.Now().UnixMilli()
		}
		_, _ = res.Array("posts")
		for _, post := range getPosts(limit, older) {
			p := post.serialize()
			_, _ = p.Array("attachments")
			for _, att := range getAttachments(post.ID) {
				_ = p.ArrayAppend(att.serialize(), "attachments")
			}
			_ = res.ArrayAppend(p, "posts")
		}
		return c.SendString(res.String())
	})

	app.Get("/post", func(c *fiber.Ctx) error {
		res := gabs.New()
		id := c.Query("id", "")
		if id == "" {
			_, _ = res.Set("ERR_INVALID_POST_ID", "error")
			return c.SendString(res.String())
		}
		post := getPost(id)
		if post == nil {
			_, _ = res.Set("ERR_UNKNOWN_POST", "error")
			return c.SendString(res.String())
		}
		res = post.serialize()
		_, _ = res.Set(post.Content, "content")
		_, _ = res.Array("attachments")
		for _, att := range getAttachments(id) {
			_ = res.ArrayAppend(att.serialize(), "attachments")
		}
		return c.SendString(res.String())
	})

	app.Post("/change-post", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		user := getProfile(txt)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		if !user.Admin {
			_, _ = res.Set("ERR_NO_PERMISSION", "error")
			return c.SendString(res.String())
		}

		payload := struct {
			Id                string   `json:"id,omitempty"`
			Title             string   `json:"title,omitempty"`
			Content           string   `json:"content,omitempty"`
			RemoveAttachments []string `json:"remove_attachments,omitempty"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
			return c.SendString(res.String())
		}

		if len(payload.Title) < 5 {
			_, _ = res.Set("ERR_POST_TITLE_MIN", "error")
			return c.SendString(res.String())
		} else if len(payload.Title) > 300 {
			_, _ = res.Set("ERR_POST_TITLE_MAX", "error")
			return c.SendString(res.String())
		}

		if len(payload.Content) < 10 {
			_, _ = res.Set("ERR_POST_CONTENT_MIN", "error")
			return c.SendString(res.String())
		} else if len(payload.Content) > 100000 {
			_, _ = res.Set("ERR_POST_CONTENT_MAX", "error")
			return c.SendString(res.String())
		}

		if len(payload.RemoveAttachments) > 0 {
			if !removeAttachments(payload.Id, payload.RemoveAttachments) {
				_, _ = res.Set("ERR_ATTACHMENT_REMOVE_FAILED", "error")
				return c.SendString(res.String())
			}
		}
		p := postChange(payload.Id, payload.Title, payload.Content)
		_, _ = res.Set(true, "success")
		_, _ = res.Set(p.ID, "id")
		return c.SendString(res.String())
	})

	app.Post("/remove-post", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		user := getProfile(txt)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		if !user.Admin {
			_, _ = res.Set("ERR_NO_PERMISSION", "error")
			return c.SendString(res.String())
		}

		id := c.Get("id")
		_, _ = res.Set(removePost(id), "success")
		return c.SendString(res.String())
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		requester := getProfile(txt)
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
	})

	app.Post("/change-users", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		requester := getProfile(txt)
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

		_, _ = res.Set(changeUserCert(payload.Certified, class), "success")

		if requester.Admin {
			_, _ = res.Set(changeUserMod(payload.Mod), "success")
		}

		return c.SendString(res.String())
	})

	app.Post("/upload-attachment", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		user := getProfile(txt)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		if !user.Admin {
			_, _ = res.Set("ERR_NO_PERMISSION", "error")
			return c.SendString(res.String())
		}

		id := c.Get("id")
		post := getPost(id)
		if post == nil {
			_, _ = res.Set("ERR_UNKNOWN_POST", "error")
			return c.SendString(res.String())
		}

		t := c.Get("content-type")

		if t == "image/png" {
			_, _ = res.Set(uploadAttachment(id, c.Body(), ".png"), "success")
		} else if t == "image/jpeg" {
			_, _ = res.Set(uploadAttachment(id, c.Body(), ".jpeg"), "success")
		} else {
			_, _ = res.Set("ERR_ILLEGAL_ATTACHMENT", "error")
			return c.SendString(res.String())
		}

		return c.SendString(res.String())
	})

	app.Post("/get-user-stats", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		requester := getProfile(txt)
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
	})

	app.Get("/events", func(c *fiber.Ctx) error {
		res := gabs.New()
		limit, err1 := strconv.Atoi(c.Query("limit", ""))
		if err1 != nil || limit > 50 {
			limit = 50
		}
		older, err2 := strconv.ParseInt(c.Query("older", ""), 10, 64)
		if err2 != nil {
			older = time.Now().UnixMilli()
		}
		fromDate, err3 := strconv.ParseInt(c.Query("from-date", ""), 10, 64)
		if err3 != nil {
			fromDate = 0
		}
		toDate, err4 := strconv.ParseInt(c.Query("to-date", ""), 10, 64)
		if err4 != nil {
			toDate = 0
		}
		_, _ = res.Array("events")
		for _, ev := range getEvents(limit, older, fromDate, toDate) {
			_ = res.ArrayAppend(ev.serialize(), "events")
		}
		return c.SendString(res.String())
	})

	app.Post("/remove-event", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		user := getProfile(txt)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		if !user.Admin {
			_, _ = res.Set("ERR_NO_PERMISSION", "error")
			return c.SendString(res.String())
		}

		id := c.Get("id")
		_, _ = res.Set(removeEvent(id), "success")
		return c.SendString(res.String())
	})

	app.Get("/event", func(c *fiber.Ctx) error {
		res := gabs.New()
		id := c.Query("id", "")
		if id == "" {
			_, _ = res.Set("ERR_INVALID_POST_ID", "error")
			return c.SendString(res.String())
		}
		event := getEvent(id)
		if event == nil {
			_, _ = res.Set("ERR_UNKNOWN_POST", "error")
			return c.SendString(res.String())
		}
		res = event.serialize()
		return c.SendString(res.String())
	})

	app.Post("/change-event", func(c *fiber.Ctx) error {
		res := gabs.New()
		token := c.Get("token")
		success, txt := analyzeTokenToEmail(token, c.UserContext())
		if !success {
			_, _ = res.Set(txt, "error")
			return c.SendString(res.String())
		}
		user := getProfile(txt)
		if user == nil {
			_, _ = res.Set("ERR_UNKNOWN_USER", "error")
			return c.SendString(res.String())
		}
		if !user.Admin {
			_, _ = res.Set("ERR_NO_PERMISSION", "error")
			return c.SendString(res.String())
		}

		payload := struct {
			Id        string `json:"id,omitempty"`
			Title     string `json:"title,omitempty"`
			StartDate int64  `json:"start_date,omitempty"`
			EndDate   int64  `json:"end_date,omitempty"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
			return c.SendString(res.String())
		}

		if len(payload.Title) < 5 {
			_, _ = res.Set("ERR_POST_TITLE_MIN", "error")
			return c.SendString(res.String())
		} else if len(payload.Title) > 300 {
			_, _ = res.Set("ERR_POST_TITLE_MAX", "error")
			return c.SendString(res.String())
		}

		if payload.StartDate > payload.EndDate {
			_, _ = res.Set("ERR_DATE_RANGE", "error")
			return c.SendString(res.String())
		}

		p := editOrCreateEvent(payload.Id, payload.Title, payload.StartDate, payload.EndDate)
		_, _ = res.Set(true, "success")
		_, _ = res.Set(p.ID, "id")
		return c.SendString(res.String())
	})

	app.Static("/static/", "./public")

	err := app.Listen(":3002")
	if err != nil {
		return
	}
}
