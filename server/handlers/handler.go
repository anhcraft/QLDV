package handlers

import (
	"das/models"
	"das/security"
	"das/utils"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"strings"
)

var guestUser = &models.User{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: security.RoleGuest,
}

var rootUser = &models.User{
	ID:   0, // pick 0 since MySQL auto-increment counter starts from 1
	Role: security.RoleRoot,
}

// GetRequester Returns the user who did the request
// - 1st return value is always not-null; if the authentication failed, it returns a "dummy" guest
// - 2nd return value is the error code; will be empty in case of success
func GetRequester(c *fiber.Ctx) (*models.User, string) {
	// TODO DON'T CACHE ACCESS TOKEN (SECURITY CONCERNS)
	token := strings.TrimSpace(c.Get("access-token"))
	if token != "" {
		res, txt := security.GetEmailFromToken(c.UserContext(), token)
		if res {
			user := getUserByEmail(txt)
			if user == nil {
				return guestUser, utils.ErrUnknownUser
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
