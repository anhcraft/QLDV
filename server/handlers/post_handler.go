package handlers

import (
	"das/models"
	"das/storage"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func getPost(id interface{}) *models.Post {
	var post models.Post
	result := storage.Db.Take(&post, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getPost while processing DB transaction")
		return nil
	} else {
		return &post
	}
}

func removePost(id interface{}) bool {
	var post models.Post
	tx := storage.Db.Where("id = ?", id).Delete(&post)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #removePost while processing DB transaction")
		return false
	}
	return true
}

func PostGetRouteHandler(c *fiber.Ctx) error {

}
