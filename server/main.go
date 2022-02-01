package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
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
	err = db_.AutoMigrate(&User{}, &Rate{}, &Achievement{}, &Post{}, &Attachment{}, &Event{}, &PostStat{})
	if err != nil {
		log.Fatalf("error migrating: %v\n", err)
	}
	db = db_
}

func getEmailFromToken(token string, c context.Context) (bool, string) {
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

	app.Post("/profile", profileGetRouteHandler)
	app.Post("/progression", progressionGetRouteHandler)
	app.Post("/change-progression", progressionChangeRouteHandler)
	app.Post("/users", userListRouteHandler)
	app.Post("/change-users", userChangeRouteHandler)
	app.Post("/get-user-stats", userStatGetRouteHandler)
	app.Post("/set-profile-cover", profileCoverSetRouteHandler)

	app.Get("/posts", postListRouteHandler)
	app.Get("/post", postGetRouteHandler)
	app.Post("/change-post", postChangeRouteHandler)
	app.Post("/remove-post", postRemoveRouteHandler)
	app.Post("/update-post-stat", postStatUpdateRouteHandler)
	app.Post("/upload-attachment", attachmentUploadRouteHandler)

	app.Get("/event", eventGetRouteHandler)
	app.Get("/events", eventListRouteHandler)
	app.Post("/remove-event", eventRemoveRouteHandler)
	app.Post("/change-event", eventChangeRouteHandler)

	app.Static("/static/", "./public")

	err := app.Listen(":3002")
	if err != nil {
		return
	}
}
