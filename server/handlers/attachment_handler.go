package handlers

import (
	"crypto/sha256"
	"das"
	"das/models"
	"encoding/hex"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	"os"
	"strconv"
	"time"
)

func uploadAttachment(postId int, data []byte, ext string) bool {
	_ = os.Mkdir("public", os.ModePerm)
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(postId) + time.Now().String()))
	md := hash.Sum(nil)
	attId := hex.EncodeToString(md)
	att := models.Attachment{
		ID:     attId + ext,
		PostId: postId,
		Date:   time.Now().UnixMilli(),
	}
	err := os.WriteFile("./public/"+attId+ext, data, os.ModePerm)
	if err != nil {
		return false
	}
	_ = main.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&att)
	return true
}

func getAttachments(postId int) []models.Attachment {
	var atts []models.Attachment
	_ = main.Db.Where("post_id = ?", postId).Order("date desc").Find(&atts)
	return atts
}

func removeAttachments(postId int, atts []string) bool {
	var att models.Attachment
	for _, v := range atts {
		main.Db.Where("id = ?", v).Where("post_id = ?", postId).Delete(&att)
	}
	return true
}

func AttachmentUploadRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := main.GetEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(emailOrError, "error")
		return c.SendString(res.String())
	}
	user := getUserByEmail(emailOrError)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	id, err := strconv.Atoi(c.Get("id"))
	if err != nil {
		_, _ = res.Set("ERR_INVALID_POST_ID", "error")
		return c.SendString(res.String())
	}
	post := main.getPost(id)
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
