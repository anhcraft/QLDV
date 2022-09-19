package handlers

import (
	"crypto/sha256"
	"das/models"
	"das/models/request"
	"das/security"
	"das/storage"
	"das/utils"
	"encoding/hex"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strconv"
	"strings"
)

const MaxProfileBoardLength = 10000
const MinProfileBoardLength = 10
const UserListLimit = 50
const MaxProfileCoverSize = 500000 // 500KB

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

func getUserById(id interface{}) *models.User {
	var user models.User
	result := storage.Db.Take(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getUserById while processing DB transaction")
		return nil
	} else {
		return &user
	}
}

func getAchievementById(id interface{}) []models.Achievement {
	var achievements []models.Achievement
	result := storage.Db.Find(&achievements, "user_id = ?", id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getAchievementById while processing DB transaction")
	}
	return achievements
}

func getAnnualRankById(id interface{}) []models.AnnualRank {
	var ranks []models.AnnualRank
	result := storage.Db.Find(&ranks, "user_id = ?", id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getAnnualRankById while processing DB transaction")
	}
	return ranks
}

func setProfileBoard(id interface{}, text string) bool {
	tx := storage.Db.Model(&models.User{}).Where("user_id = ?", id).Update("profile_board", text)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setProfileBoard while processing DB transaction")
		return false
	}
	return tx.RowsAffected > 0
}

func setProfileSettings(id interface{}, settings uint8) bool {
	tx := storage.Db.Model(&models.User{}).Where("user_id = ?", id).Update("profile_settings", settings)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setProfileSettings while processing DB transaction")
		return false
	}
	return tx.RowsAffected > 0
}

func setRole(id interface{}, role uint8) bool {
	tx := storage.Db.Model(&models.User{}).Where("user_id = ?", id).Update("role", role)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setRole while processing DB transaction")
		return false
	}
	return tx.RowsAffected > 0
}

func setAchievements(id interface{}, achievements []models.Achievement) bool {
	var ach models.Achievement
	tx := storage.Db.Where("user_id = ?", id).Delete(&ach)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setAchievements while processing DB transaction (1)")
		return false
	}
	tx = storage.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(achievements)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setAchievements while processing DB transaction (2)")
		return false
	}
	return tx.RowsAffected > 0
}

func setAnnualRanks(id interface{}, annualRanks []models.AnnualRank) bool {
	var ar models.AnnualRank
	tx := storage.Db.Where("user_id = ?", id).Delete(&ar)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setAnnualRanks while processing DB transaction (1)")
		return false
	}
	tx = storage.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(annualRanks)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setAnnualRanks while processing DB transaction (2)")
		return false
	}
	return tx.RowsAffected > 0
}

func getUsers(req *request.UserListModel) []models.User {
	var users []models.User
	cmd := storage.Db.Limit(int(req.Limit))
	if req.BelowId > 0 {
		cmd = cmd.Where("id < ?", req.BelowId)
	}
	if len(req.FilterName) > 0 {
		cmd = cmd.Where("LOWER(`name`) like ?", "%"+req.FilterName+"%")
	}
	if len(req.FilterClass) > 0 {
		cmd = cmd.Where("LOWER(`class`) like ?", "%"+req.FilterClass+"%")
	}
	if len(req.FilterEmail) > 0 {
		cmd = cmd.Where("LOWER(`email`) like ?", "%"+req.FilterEmail+"%")
	}
	if utils.IsLoggedIn(req.FilterRole) {
		cmd = cmd.Where("role = ?", req.FilterRole)
	}
	cmd = cmd.Find(&users)
	if cmd.Error != nil {
		log.Error().Err(cmd.Error).Msg("An error occurred at #getUsers while processing DB transaction")
	}
	return users
}

func setProfileCover(id uint16, data []byte, ext string) (bool, string) {
	_ = os.Mkdir("public", os.ModePerm)
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(int(id))))
	fileName := "cover-" + hex.EncodeToString(hash.Sum(nil)) + ext
	err := os.WriteFile("./public/"+fileName, data, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred at #setProfileCover while writing file")
		return false, ""
	}
	tx := storage.Db.Model(&models.User{}).Where("user_id = ?", id).Update("profile_cover", fileName)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setProfileCover while processing DB transaction")
		return false, ""
	}
	return tx.RowsAffected > 0, fileName
}

func UserGetRouteHandler(c *fiber.Ctx) error {
	queryTable := struct {
		profile      bool `query:"profile"`
		achievements bool `query:"achievements"`
		annualRanks  bool `query:"annual-ranks"`
	}{}

	if err := c.QueryParser(&queryTable); err != nil {
		log.Error().Err(err).Msg("There was an error occurred while parsing queries at #UserGetRouteHandler")
		return ReturnError(c, utils.ErrInvalidRequestQuery)
	}

	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	whoParam := c.Params("id")
	var who *models.User
	if whoParam == "" {
		who = requester
	} else {
		if !utils.ValidateNonNegativeInteger(whoParam) {
			return ReturnError(c, utils.ErrInvalidToken)
		}
		who = getUserById(whoParam)
	}
	if who == nil || who.Role == utils.RoleGuest {
		return ReturnError(c, utils.ErrUnknownUser)
	}

	data := gabs.New()
	if queryTable.profile {
		_, _ = data.Set(who.Serialize(requester), "profile")
	}
	if queryTable.achievements && (requester.HasPrivilegeOver(who, 0) || who.IsAchievementPublic()) {
		_, _ = data.Array("achievements")
		for _, v := range getAchievementById(who.ID) {
			_ = data.ArrayAppend(v.Serialize(), "achievements")
		}
	}
	if queryTable.annualRanks && (requester.HasPrivilegeOver(who, 0) || who.IsAnnualRankPublic()) {
		_, _ = data.Array("annualRanks")
		for _, v := range getAnnualRankById(who.ID) {
			_ = data.ArrayAppend(v.Serialize(), "annualRanks")
		}
	}
	return ReturnJSON(c, data)
}

func UserUpdateRouteHandler(c *fiber.Ctx) error {
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	whoParam := c.Params("id")
	var who *models.User
	if whoParam == "" {
		who = requester
	} else {
		if !utils.ValidateNonNegativeInteger(whoParam) {
			return ReturnError(c, utils.ErrInvalidToken)
		}
		who = getUserById(whoParam)
	}
	if who == nil || who.Role == utils.RoleGuest {
		return ReturnError(c, utils.ErrUnknownUser)
	}
	if !requester.HasPrivilegeOver(who, 0) {
		return ReturnError(c, utils.ErrNoPermission)
	}

	json, err2 := gabs.ParseJSON(c.Body())
	if err2 != nil {
		return ReturnError(c, utils.ErrInvalidRequestBody)
	}
	if requester.ID == who.ID {
		response := gabs.New()

		profileSettingDirty := false
		profileLocked := json.Path("profile.settings.profileLocked")
		if profileLocked != nil {
			profileSettingDirty = true
			if profileLocked.Data().(bool) {
				who.ProfileSettings = who.ProfileSettings | 1
			} else {
				who.ProfileSettings = who.ProfileSettings &^ 1
			}
		}
		classPublic := json.Path("profile.settings.classPublic")
		if classPublic != nil {
			profileSettingDirty = true
			if classPublic.Data().(bool) {
				who.ProfileSettings = who.ProfileSettings | 2
			} else {
				who.ProfileSettings = who.ProfileSettings &^ 2
			}
		}
		achievementPublic := json.Path("profile.settings.achievementPublic")
		if achievementPublic != nil {
			profileSettingDirty = true
			if achievementPublic.Data().(bool) {
				who.ProfileSettings = who.ProfileSettings | 4
			} else {
				who.ProfileSettings = who.ProfileSettings &^ 4
			}
		}
		annualRankPublic := json.Path("profile.settings.annualRankPublic")
		if annualRankPublic != nil {
			profileSettingDirty = true
			if annualRankPublic.Data().(bool) {
				who.ProfileSettings = who.ProfileSettings | 8
			} else {
				who.ProfileSettings = who.ProfileSettings &^ 8
			}
		}
		if profileSettingDirty {
			_, _ = response.Set(setProfileSettings(who.ID, who.ProfileSettings), "profile.settings")
		}

		profileBoard := json.Path("profile.profileBoard")
		if profileBoard != nil {
			who.ProfileBoard = profileBoard.Data().(string)
			who.ProfileBoard = strings.TrimSpace(who.ProfileBoard)
			who.ProfileBoard = security.SafeHTMLPolicy.Sanitize(who.ProfileBoard)

			if len(who.ProfileBoard) < MinProfileBoardLength {
				return ReturnError(c, utils.ErrProfileBoardTooShort)
			} else if len(who.ProfileBoard) > MaxProfileBoardLength {
				return ReturnError(c, utils.ErrProfileBoardTooLong)
			}

			_, _ = response.Set(setProfileBoard(who.ID, who.ProfileBoard), "profile.profileBoard")
		}

		return ReturnJSON(c, response)

	} else if utils.IsManager(requester.Role) {
		isClassScope := utils.GetRoleGroup(requester.Role) == utils.RoleGroupClassManager
		response := gabs.New()

		if json.Exists("profile.role") {
			role := json.Path("profile.role").Data().(uint8)
			if isClassScope {
				if who.Class != requester.Class ||
					utils.GetRoleGroup(who.Role) != utils.RoleGroupMember ||
					utils.GetRoleGroup(role) != utils.RoleGroupMember {
					return ReturnError(c, utils.ErrNoPermission)
				}
			} else {
				if (utils.GetRoleGroup(who.Role) != utils.RoleGroupMember &&
					utils.GetRoleGroup(who.Role) != utils.RoleGroupClassManager) ||
					(utils.GetRoleGroup(role) != utils.RoleGroupMember &&
						utils.GetRoleGroup(role) != utils.RoleGroupClassManager) {
					return ReturnError(c, utils.ErrNoPermission)
				}
			}
			_, _ = response.Set(setRole(who.ID, role), "profile.role")
		}

		if json.Exists("achievements") {
			ach := make([]models.Achievement, 0)
			for _, child := range json.Path("achievements").Children() {
				title := child.Path("title").Data().(string)
				title = strings.TrimSpace(title)
				title = bluemonday.StripTagsPolicy().Sanitize(title)
				if len(title) == 0 {
					continue
				}
				year := child.Path("year").Data().(uint16)
				if year < who.EntryYear || year > who.EntryYear+2 { // N, N+1, N+2
					continue
				}
				ach = append(ach, models.Achievement{
					UserId: who.ID,
					Title:  title,
					Year:   year,
				})
			}
			_, _ = response.Set(setAchievements(who.ID, ach), "achievements")
		}

		if json.Exists("annualRanks") {
			ar := make([]models.AnnualRank, 0)
			for _, child := range json.Path("annualRanks").Children() {
				year := child.Path("year").Data().(uint16)
				if year < who.EntryYear || year > who.EntryYear+2 { // N, N+1, N+2
					continue
				}
				level := child.Path("level").Data().(uint8)
				if level < models.UnknownRank || level > models.MediumRank {
					continue
				}
				ar = append(ar, models.AnnualRank{
					UserId: who.ID,
					Level:  level,
					Year:   year,
				})
			}
			_, _ = response.Set(setAnnualRanks(who.ID, ar), "annualRanks")
		}

		return ReturnJSON(c, response)
	} else {
		return ReturnError(c, utils.ErrNoPermission)
	}
}

func UserListRouteHandler(c *fiber.Ctx) error {
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if !utils.IsManager(requester.Role) {
		return ReturnError(c, utils.ErrNoPermission)
	}
	req := request.UserListModel{}
	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("There was an error occurred while parsing body at #UserListRouteHandler")
		return ReturnError(c, utils.ErrInvalidRequestBody)
	}
	req.FilterName = strings.ToLower(utils.RemoveVietnameseAccents(strings.TrimSpace(req.FilterName)))
	req.FilterClass = strings.ToLower(utils.RemoveVietnameseAccents(strings.TrimSpace(req.FilterClass)))
	req.FilterEmail = strings.ToLower(utils.RemoveVietnameseAccents(strings.TrimSpace(req.FilterEmail)))
	if utils.GetRoleGroup(requester.Role) == utils.RoleGroupClassManager {
		req.FilterClass = requester.Class
	}
	req.Limit = utils.ClampUint8(req.Limit, 0, UserListLimit)
	users := gabs.New()
	_, _ = users.Array("users")
	for _, user := range getUsers(&req) {
		_ = users.ArrayAppend(user.Serialize(requester), "users")
	}
	return ReturnJSON(c, users)
}

func UserStatGetRouteHandler(c *fiber.Ctx) error {
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if utils.GetRoleGroup(requester.Role) != utils.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}

	result := struct {
		a int64
		b int64
		c int64
		d int64
		e int64
		f int64
		g int64
		h int64
		i int64
	}{}

	cmd := "select count(if(class like '10%', 1, null)) as a"
	cmd += ", count(if(class like '11%', 1, null)) as b"
	cmd += ", count(if(class like '12%', 1, null)) as c"
	cmd += ", count(if(role = " + strconv.Itoa(int(utils.RoleRegularMember)) + ", 1, null)) as d"
	cmd += ", count(if(role = " + strconv.Itoa(int(utils.RoleCertifiedMember)) + ", 1, null)) as e"
	cmd += ", count(if(role = " + strconv.Itoa(int(utils.RoleClassDeputySecretary)) + ", 1, null)) as f"
	cmd += ", count(if(role = " + strconv.Itoa(int(utils.RoleClassSecretary)) + ", 1, null)) as g"
	cmd += ", count(if(role = " + strconv.Itoa(int(utils.RoleDeputySecretary)) + ", 1, null)) as h"
	cmd += ", count(if(role = " + strconv.Itoa(int(utils.RoleSecretary)) + ", 1, null)) as i"
	cmd += " from users"
	_ = storage.Db.Raw(cmd).Row().Scan(&result.a, &result.b, &result.c, &result.d, &result.e)
	response := gabs.New()
	_, _ = response.Set(result.a, "user-count-by-grade", "grade-10")
	_, _ = response.Set(result.b, "user-count-by-grade", "grade-11")
	_, _ = response.Set(result.c, "user-count-by-grade", "grade-12")
	_, _ = response.Set(result.d, "user-count-by-role", "regular-member")
	_, _ = response.Set(result.e, "user-count-by-role", "certified-member")
	_, _ = response.Set(result.f, "user-count-by-role", "class-deputy-secretary")
	_, _ = response.Set(result.g, "user-count-by-role", "class-secretary")
	_, _ = response.Set(result.h, "user-count-by-role", "deputy-secretary")
	_, _ = response.Set(result.i, "user-count-by-role", "secretary")
	return ReturnJSON(c, response)
}

func ProfileCoverUploadRouteHandler(c *fiber.Ctx) error {
	if len(c.Body()) > MaxProfileCoverSize {
		return ReturnError(c, utils.ErrProfileCoverTooLarge)
	}

	requester, err := GetRequester(c)
	if err != "" || !utils.IsLoggedIn(requester.Role) {
		return ReturnError(c, err)
	}

	t := c.Get("content-type")

	ok := false
	fn := ""
	// TODO Check the file content rather than the given content-type since it is inaccurate
	if t == "image/png" {
		ok, fn = setProfileCover(requester.ID, c.Body(), ".png")
	} else if t == "image/jpeg" {
		ok, fn = setProfileCover(requester.ID, c.Body(), ".jpeg")
	} else {
		return ReturnError(c, utils.ErrUnsupportedProfileCover)
	}

	if !ok {
		return ReturnError(c, utils.ErrProfileCoverUploadFailed)
	}

	response := gabs.New()
	_, _ = response.Set(fn, "name")
	return ReturnJSON(c, response)
}
