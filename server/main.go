package main

import (
	"das/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(recover2.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	app.Get("/user/:id?", handlers.UserGetRouteHandler)
	app.Post("/user/:id?", handlers.UserUpdateRouteHandler)
	app.Get("/users/", handlers.UserListRouteHandler)
	app.Get("/user-stats/", handlers.UserStatGetRouteHandler)
	app.Post("/user-profile-cover/", handlers.ProfileCoverUploadRouteHandler)

	app.Get("/post/:id", handlers.PostGetRouteHandler)
	app.Post("/post/:id?", handlers.PostUpdateRouteHandler)
	app.Delete("/post/:id", handlers.PostDeleteRouteHandler)
	app.Get("/posts/", handlers.PostListRouteHandler)
	app.Post("/post-stat/:id", handlers.PostStatUpdateRouteHandler)
	app.Post("/post-attachment/:id", handlers.AttachmentUploadRouteHandler)
	app.Delete("/post-attachment/:id", handlers.AttachmentDeleteRouteHandler)
	app.Get("/post-hashtags/", handlers.PostHashtagListRouteHandler)

	app.Get("/event/:id", handlers.EventGetRouteHandler)
	app.Post("/event/:id?", handlers.EventUpdateRouteHandler)
	app.Delete("/event/:id", handlers.EventRemoveRouteHandler)
	app.Get("/events/", handlers.EventListRouteHandler)
	/*
		//app.Post("/contest", contestGetRouteHandler)
		app.Post("/change-contest", handlers.ContestChangeRouteHandler)
		app.Post("/remove-contest", handlers.ContestRemoveRouteHandler)
		app.Post("/contest-sessions", handlers.ContestSessionListRouteHandler)
		app.Post("/get-contest-stats", handlers.ContestStatGetRouteHandler)
		app.Post("/submit-contest-session", handlers.ContestSessionSubmitRouteHandler)
		app.Post("/join-contest-session", handlers.ContestSessionJoinRouteHandler)*/

	app.Static("/static/", "./public")
	app.Get("/status/", monitor.New())

	err := app.Listen(":3002")
	if err != nil {
		return
	}
}
