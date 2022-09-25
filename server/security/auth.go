package security

import (
	"context"
	"das/utils"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/rs/zerolog/log"
)

var client *auth.Client

func init() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while initializing Firebase (1)")
		return
	}
	client_, err := app.Auth(ctx)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while initializing Firebase (2)")
		return
	}
	client = client_
}

// GetEmailFromToken Returns email by token relevant to a user
// - 1st return value: true if success
// - 2nd return value: the email (when success), or the error code
func GetEmailFromToken(c context.Context, token string) (bool, string) {
	// TODO DON'T CACHE ACCESS TOKEN (SECURITY CONCERNS)
	tkn, err := client.VerifyIDToken(c, token)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while validating Firebase ID token")
		return false, utils.ErrTokenVerify
	}
	u, err := client.GetUser(c, tkn.UID)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while getting Firebase user")
		return false, utils.ErrUnknownTokenOwner
	}
	return true, u.Email
}
