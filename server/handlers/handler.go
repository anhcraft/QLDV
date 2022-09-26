package handlers

import (
	"das/models"
	"das/security"
	"das/storage"
	"das/utils"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

var guestRequester = &Requester{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: security.RoleGuest,
}

var rootRequester = &Requester{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: security.RoleRoot,
}

var guestUser = &models.User{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: security.RoleGuest,
}

var rootUser = &models.User{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: security.RoleRoot,
}

// Internal method for fetching user
func getUserByEmail(email string) *models.User {
	var user models.User
	result := storage.Db.Take(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getUserByEmail while processing DB transaction")
		return nil
	} else {
		return &user
	}
}

// GetRequester Returns the requester who sent the request
// - 1st return value is always not-null; if the authentication failed, it returns a (dummy) guest requester
// - 2nd return value is the error code; will be empty in case of success
func GetRequester(c *fiber.Ctx) (*Requester, string) {
	// TODO DON'T CACHE ACCESS TOKEN (SECURITY CONCERNS)
	token := strings.TrimSpace(c.Get("access-token"))
	if token == "" {
		return guestRequester, ""
	}
	has, txt := security.GetEmailFromToken(c.UserContext(), token)
	if !has {
		return guestRequester, txt
	}
	cache, ok := storage.GetCache("requesters", txt)
	if ok {
		requester := cache.(Requester)
		return &requester, ""
	} else {
		user := getUserByEmail(txt)
		if user == nil {
			storage.SetCache("requesters", txt, guestRequester, 1*time.Minute)
			return guestRequester, utils.ErrUnknownTokenOwner
		} else {
			req := Requester{
				ID:   user.ID,
				Role: user.Role,
			}
			storage.SetCache("requesters", txt, req, 5*time.Minute)
			return &req, ""
		}
	}
}

// GetRequesterUser Returns the user who sent the request
// - 1st return value is always not-null; if the authentication failed, it returns a (dummy) guest user
// - 2nd return value is the error code; will be empty in case of success
func GetRequesterUser(c *fiber.Ctx) (*models.User, string) {
	// TODO DON'T CACHE ACCESS TOKEN (SECURITY CONCERNS)
	token := strings.TrimSpace(c.Get("access-token"))
	if token != "" {
		res, txt := security.GetEmailFromToken(c.UserContext(), token)
		if res {
			user := getUserByEmail(txt)
			if user == nil {
				return guestUser, utils.ErrUnknownUser
			}
			storage.SetCache("requesters", txt, Requester{
				ID:   user.ID,
				Role: user.Role,
			}, 5*time.Minute)
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

func ReturnString(c *fiber.Ctx, data string) error {
	return c.Status(fiber.StatusOK).SendString(data)
}

func BuildResponse(container *gabs.Container) string {
	res := gabs.New()
	_, _ = res.Set(container, "result")
	_, _ = res.Set(true, "success")
	return res.String()
}

func ReturnJSON(c *fiber.Ctx, container *gabs.Container) error {
	return ReturnString(c, BuildResponse(container))
}

func ReturnEmpty(c *fiber.Ctx) error {
	res := gabs.New()
	_, _ = res.Set(true, "success")
	return ReturnString(c, res.String())
}
