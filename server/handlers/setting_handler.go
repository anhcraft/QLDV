package handlers

import (
	"das/models"
	"das/models/request"
	"das/storage"
	"das/utils"
	"encoding/json"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

const HomepageSettingKey = "homepage"

func getSetting(key string, requester *models.User) string {
	cache, ok := storage.GetCache("settings", key)
	if ok {
		data := cache.(models.Settings)
		if data.Privacy > requester.Role {
			return "{}"
		}
		return data.Value
	}

	var settings models.Settings
	result := storage.Db.Take(&settings, "`key` = ? and privacy <= ?", key, requester.Role)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "{}"
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getSetting while processing DB transaction")
		return "{}"
	} else {
		storage.SetCache("settings", key, settings, 5*time.Minute)
		return settings.Value
	}
}

func updateSetting(key string, val string, requester *models.User) bool {
	tx := storage.Db.Model(&models.Settings{}).Where("`key` = ? and privacy <= ?", key, requester.Role).Update("value", val)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #updateSetting while processing DB transaction")
	}
	if tx.RowsAffected > 0 {
		storage.PurgeCache("settings", key)
		return true
	}
	return false
}

func SettingGetRouteHandler(c *fiber.Ctx) error {
	requester, err1 := GetRequester(c)
	if err1 != "" {
		return ReturnError(c, err1)
	}
	param := c.Params("id")
	res, err2 := gabs.ParseJSON([]byte(getSetting(param, requester)))
	if err2 != nil {
		log.Error().Err(err2).Msg("An error occurred at #SettingGetRouteHandler while parsing settings data")
		return ReturnEmpty(c)
	}
	return ReturnJSON(c, res)
}

func SettingUpdateRouteHandler(c *fiber.Ctx) error {
	requester, err1 := GetRequester(c)
	if err1 != "" {
		return ReturnError(c, err1)
	}
	param := c.Params("id")
	data := "{}"

	// Validation
	if param == HomepageSettingKey {
		req := &request.HomepageSettingModel{}
		if err2 := c.BodyParser(&req); err2 != nil {
			log.Error().Err(err2).Str("key", HomepageSettingKey).Msg("There was an error occurred while parsing body at #SettingUpdateRouteHandler")
			return ReturnError(c, utils.ErrInvalidRequestBody)
		}
		req.FeaturedUserLimit = utils.ClampUint8(req.FeaturedUserLimit, 0, 10)
		req.FeaturedAchievementLimit = utils.ClampUint8(req.FeaturedAchievementLimit, 0, 10)
		req.ActivitySlideshow = ClearEmpty(req.ActivitySlideshow)
		_data, err3 := json.Marshal(req)
		if err3 != nil {
			log.Error().Err(err3).Str("key", HomepageSettingKey).Msg("There was an error occurred while serializing data at #SettingUpdateRouteHandler")
			return ReturnError(c, utils.ErrSettingUpdateFailed)
		}
		data = string(_data)
	}

	if updateSetting(param, data, requester) {
		return ReturnEmpty(c)
	} else {
		return ReturnError(c, utils.ErrSettingUpdateFailed)
	}
}

func ClearEmpty(list []string) []string {
	q := make([]string, 0)
	for _, v := range list {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			q = append(q, v)
		}
	}
	return q
}
