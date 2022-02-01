package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	"os"
	"time"
)

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

func attachmentUploadRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	user := getProfile(email)
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
}
