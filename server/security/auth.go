package security

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"log"
	"strings"
)

var client *auth.Client

func init() {
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

func GetEmailFromToken(token string, c context.Context) (bool, string) {
	token = strings.TrimSpace(token)
	if len(token) == 0 {
		return false, "ERR_TOKEN_VERIFY"
	}
	tkn, err := client.VerifyIDToken(c, token)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return false, "ERR_TOKEN_VERIFY_BACKEND"
	}
	u, err := client.GetUser(c, tkn.UID)
	if err != nil {
		log.Printf("error getting user %s: %v\n", tkn.UID, err)
		return false, "ERR_USER_GET_BACKEND"
	}
	return true, u.Email
}
