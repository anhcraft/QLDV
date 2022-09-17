package handlers

import (
	"context"
	"das/models"
	"das/utils"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

var guestUser = &models.User{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: utils.RoleGuest,
}
var client *auth.Client

func setupFirebase() bool {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		_ = fmt.Errorf("error initializing app: %v", err)
		return false
	}
	client_, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return false
	}
	client = client_
	return true
}

// GetEmailFromToken Returns email by token relevant to a user
// - 1st return value: true if success
// - 2nd return value: the email (when success), or the error code
func GetEmailFromToken(c context.Context, token string) (bool, string) {
	tkn, err := client.VerifyIDToken(c, token)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while validating Firebase ID token")
		return false, ErrTokenVerify
	}
	u, err := client.GetUser(c, tkn.UID)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while getting Firebase user")
		return false, ErrTokenUnknownUser
	}
	return true, u.Email
}

// GetRequester Returns the user who did the request
// - 1st return value is always not-null; if the authentication failed, it returns a "dummy" guest
// - 2nd return value is the error code; will be empty in case of success
func GetRequester(c *fiber.Ctx) (*models.User, string) {
	token := strings.TrimSpace(c.Get("access-token"))
	if token != "" {
		res, txt := GetEmailFromToken(c.UserContext(), token)
		if res {
			user := getUserByEmail(txt)
			if user == nil {
				return guestUser, ErrUserNotExist
			}
			return user, ""
		}
		return guestUser, txt
	}
	return guestUser, ""
}

func ReturnError(c *fiber.Ctx, err string) error {
	res := gabs.New()
	_, _ = res.Set(err, "error")
	_, _ = res.Set(false, "success")
	return c.Status(fiber.StatusBadRequest).SendString(res.String())
}

func ReturnJSON(c *fiber.Ctx, container *gabs.Container) error {
	res := gabs.New()
	_, _ = res.Set(container, "result")
	_, _ = res.Set(true, "success")
	return c.Status(fiber.StatusOK).SendString(res.String())
}
