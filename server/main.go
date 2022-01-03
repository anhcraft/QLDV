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
	"strconv"
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
	dsn := "sql6462579:rgxWwnqF21@tcp(sql6.freemysqlhosting.net:3306)/sql6462579?charset=utf8mb4&parseTime=True&loc=Local"
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting database: %v\n", err)
	}
	err = db_.AutoMigrate(&User{}, &Rate{}, &Achievement{}, &Post{})
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

func getPost(id string) *Post {
	var post Post
	result := db.Take(&post, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &post
	}
}

func postChange(id string, title string, content string) bool {
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
	return true
}

func getPosts(limit int, older int64) []Post {
	var posts []Post
	_ = db.Where("date < ?", older).Order("date desc").Limit(limit).Find(&posts)
	return posts
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
		_, _ = res.Array("rates")
		for _, rate := range getRates(txt) {
			_ = res.ArrayAppend(rate.serialize(), "rates")
		}
		_, _ = res.Array("achievements")
		for _, rate := range getAchievements(txt) {
			_ = res.ArrayAppend(rate.serialize(), "achievements")
		}
		return c.SendString(res.String())
	})

	app.Get("/posts", func(c *fiber.Ctx) error {
		res := gabs.New()
		limit, err1 := strconv.Atoi(c.Query("limit", ""))
		if err1 != nil || limit > 10 {
			limit = 10
		}
		older, err2 := strconv.ParseInt(c.Query("older", ""), 10, 64)
		if err2 != nil {
			older = time.Now().UnixMilli()
		}
		_, _ = res.Array("posts")
		for _, post := range getPosts(limit, older) {
			_ = res.ArrayAppend(post.serialize(), "posts")
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
			Id      string `json:"id,omitempty"`
			Title   string `json:"title,omitempty"`
			Content string `json:"content,omitempty"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
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

		_, _ = res.Set(postChange(payload.Id, payload.Title, payload.Content), "success")
		return c.SendString(res.String())
	})

	err := app.Listen(":3002")
	if err != nil {
		return
	}
}
