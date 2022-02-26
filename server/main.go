package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"regexp"
)

var client *auth.Client
var db *gorm.DB
var ugcPolicy *bluemonday.Policy

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
	err = db_.AutoMigrate(&User{}, &Rate{}, &Achievement{}, &Post{}, &Attachment{}, &Event{}, &PostStat{}, &Contest{})
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
	ugcPolicy = bluemonday.UGCPolicy()
	ugcPolicy.AllowStyles("color", "background-color", "text-align", "width", "height", "font-size", "font-weight", "padding-left").Globally()
	ugcPolicy.AddTargetBlankToFullyQualifiedLinks(true)
	ugcPolicy.AllowElements("iframe")
	ugcPolicy.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	ugcPolicy.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	ugcPolicy.AllowAttrs("src").OnElements("iframe")
	ugcPolicy.AllowAttrs("frameborder").Matching(bluemonday.Number).OnElements("iframe")
	ugcPolicy.AllowAttrs("allow").Matching(regexp.MustCompile(`[a-z; -]*`)).OnElements("iframe")
	ugcPolicy.AllowAttrs("allowfullscreen").OnElements("iframe")

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
	app.Post("/set-profile-board", profileBoardSetRouteHandler)

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
	//app.Post("/contest", contestGetRouteHandler)
	app.Post("/change-contest", contestChangeRouteHandler)
	app.Post("/remove-contest", contestRemoveRouteHandler)

	app.Static("/static/", "./public")

	err := app.Listen(":3002")
	if err != nil {
		return
	}
}
