package handlers

import (
	"crypto/sha256"
	"das/models"
	"das/security"
	"das/storage"
	"das/utils"
	"encoding/hex"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
	"os"
	"strconv"
	"time"
)

const MaxAttachmentSize = 2000000 // 2MB

func uploadAttachment(postId uint32, data []byte, ext string) (bool, string, string) {
	_ = os.Mkdir("public", os.ModePerm)
	hash := sha256.New()
	hash.Write([]byte(strconv.FormatUint(uint64(postId), 10) + time.Now().String() + ext))
	id := hex.EncodeToString(hash.Sum(nil))
	fileName := id + ext
	err := os.WriteFile("./public/"+fileName, data, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred at #uploadAttachment while writing file")
		return false, "", ""
	}
	att := models.Attachment{
		ID:     id,
		PostId: postId,
	}
	tx := storage.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&att)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #uploadAttachment while processing DB transaction")
		return false, "", ""
	}
	return tx.RowsAffected > 0, id, fileName
}

func getAttachments(postId uint32) []models.Attachment {
	var atts []models.Attachment
	cmd := storage.Db.Where("post_id = ?", postId).Order("updateDate desc").Find(&atts)
	if cmd.Error != nil {
		log.Error().Err(cmd.Error).Msg("An error occurred at #getAttachments while processing DB transaction")
	}
	return atts
}

func removeAttachment(attId string, privacy uint8) bool {
	var att models.Attachment
	cmd := storage.Db.Joins("right join posts on attachments.post_id = posts.id").Where("attachments.id = ?", attId).Where("posts.privacy <= ?", privacy).Delete(&att)
	if cmd.Error != nil {
		log.Error().Err(cmd.Error).Msg("An error occurred at #removeAttachment while processing DB transaction")
	}
	return cmd.RowsAffected > 0
}

func AttachmentUploadRouteHandler(c *fiber.Ctx) error {
	if len(c.Body()) > MaxAttachmentSize {
		return ReturnError(c, utils.ErrAttachmentTooLarge)
	}

	id := c.Params("id")
	if id == "" || !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if security.GetRoleGroup(requester.Role) != security.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}

	post := getPost(id)
	if post == nil {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	if post.Privacy > requester.Role {
		return ReturnError(c, utils.ErrNoPermission)
	}

	t := c.Get("content-type")

	ok := false
	attId := ""
	fn := ""
	// TODO Check the file content rather than the given content-type since it is inaccurate
	if t == "image/png" {
		ok, attId, fn = uploadAttachment(post.ID, c.Body(), ".png")
	} else if t == "image/jpeg" {
		ok, attId, fn = uploadAttachment(post.ID, c.Body(), ".jpeg")
	} else {
		return ReturnError(c, utils.ErrUnsupportedAttachment)
	}

	if !ok {
		return ReturnError(c, utils.ErrAttachmentUploadFailed)
	}

	response := gabs.New()
	_, _ = response.Set(fn, "name")
	_, _ = response.Set(attId, "id")
	return ReturnJSON(c, response)
}

func AttachmentDeleteRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id") // attachment ID
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if security.GetRoleGroup(requester.Role) != security.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}

	if removeAttachment(id, requester.Role) {
		return ReturnEmpty(c)
	} else {
		return ReturnError(c, utils.ErrAttachmentDeleteFailed)
	}
}
