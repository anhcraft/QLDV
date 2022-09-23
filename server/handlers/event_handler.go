package handlers

import (
	"das/models"
	"das/models/request"
	"das/security"
	"das/storage"
	"das/utils"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const MaxEventTitleLength = 300
const MinEventTitleLength = 10
const EventListLimit = 50

func getEvents(req *request.EventListModel, requester *models.User) []models.Event {
	var events []models.Event
	a := storage.Db.Limit(int(req.Limit)).Order("end_date desc, begin_date desc, id desc")
	a = a.Where("privacy <= ?", requester.Role)
	if req.BeginDate > 0 && req.EndDate > 0 && req.BeginDate <= req.EndDate {
		a = a.Where("begin_date <= ? and end_date >= ?", req.EndDate, req.BeginDate)
	} else if req.BeginDate > 0 && req.EndDate == 0 {
		a = a.Where("begin_date > ?", req.BeginDate)
	} else if req.BeginDate == 0 && req.EndDate > 0 {
		a = a.Where("end_date < ?", req.EndDate)
	}
	if req.BelowId > 0 {
		a = a.Where("id < ?", req.BeginDate)
	}
	a = a.Find(&events)
	if a.Error != nil {
		log.Error().Err(a.Error).Msg("An error occurred at #getEvents while processing DB transaction")
	}
	return events
}

func removeEvent(id interface{}) bool {
	var event models.Event
	tx := storage.Db.Where("id = ?", id).Delete(&event)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #removeEvent while processing DB transaction")
		return false
	}
	return tx.RowsAffected > 0
}

func getEvent(id interface{}) *models.Event {
	var event models.Event
	result := storage.Db.Take(&event, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getEvent while processing DB transaction")
		return nil
	} else {
		return &event
	}
}

func updateOrCreateEvent(id uint32, req *request.EventUpdateModel) *models.Event {
	event := models.Event{
		Link:      utils.GenerateLinkFromTitle(req.Title),
		Title:     req.Title,
		BeginDate: req.BeginDate,
		EndDate:   req.EndDate,
		Privacy:   req.Privacy,
	}
	if id > 0 {
		event.ID = id
	}
	tx := storage.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "link", "begin_date", "end_date", "privacy"}),
	}).Create(&event)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #updateOrCreateEvent while processing DB transaction")
	} else if tx.RowsAffected > 0 {
		return &event
	}
	return nil
}

func EventGetRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" || !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownEvent)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}

	event := getEvent(id)
	if event == nil {
		return ReturnError(c, utils.ErrUnknownEvent)
	}
	if event.Privacy > requester.Role {
		return ReturnError(c, utils.ErrNoPermission)
	}
	return ReturnJSON(c, event.Serialize())
}

func EventRemoveRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" || !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownEvent)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if security.GetRoleGroup(requester.Role) < security.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}

	event := getEvent(id)
	if event == nil {
		return ReturnError(c, utils.ErrUnknownEvent)
	}
	if event.Privacy > requester.Role {
		return ReturnError(c, utils.ErrNoPermission)
	}
	if removeEvent(event.ID) {
		return ReturnEmpty(c)
	} else {
		return ReturnError(c, utils.ErrEventDeleteFailed)
	}
}

func EventUpdateRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id != "" && !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownEvent)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if security.GetRoleGroup(requester.Role) < security.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}
	eventId := uint32(0)
	if id != "" {
		event := getEvent(id)
		if event == nil {
			return ReturnError(c, utils.ErrUnknownEvent)
		}
		if event.Privacy > requester.Role {
			return ReturnError(c, utils.ErrNoPermission)
		}
		eventId = event.ID
	}

	req := &request.EventUpdateModel{}
	if err2 := c.BodyParser(&req); err2 != nil {
		log.Error().Err(err2).Msg("There was an error occurred while parsing body at #EventUpdateRouteHandler")
		return ReturnError(c, utils.ErrInvalidRequestBody)
	}
	req.Title = bluemonday.StrictPolicy().Sanitize(strings.TrimSpace(req.Title))
	if req.Privacy < 0 {
		req.Privacy = 0
	}

	if len(req.Title) < MinEventTitleLength {
		return ReturnError(c, utils.ErrEventTitleTooShort)
	} else if len(req.Title) > MaxEventTitleLength {
		return ReturnError(c, utils.ErrEventTitleTooLong)
	}

	if req.EndDate < req.BeginDate || req.BeginDate < 0 || req.EndDate < 0 {
		return ReturnError(c, utils.ErrEventInvalidDuration)
	}

	event := updateOrCreateEvent(eventId, req)
	if event == nil {
		if eventId == 0 {
			return ReturnError(c, utils.ErrEventCreateFailed)
		} else {
			return ReturnError(c, utils.ErrEventUpdateFailed)
		}
	}
	response := gabs.New()
	_, _ = response.Set(event.ID, "id")
	return ReturnJSON(c, response)
}

func EventListRouteHandler(c *fiber.Ctx) error {
	req := request.EventListModel{}
	if err := c.QueryParser(&req); err != nil {
		log.Error().Err(err).Msg("There was an error occurred while parsing queries at #EventListRouteHandler")
		return ReturnError(c, utils.ErrInvalidRequestQuery)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	req.Limit = utils.ClampUint8(req.Limit, 0, EventListLimit)

	events := gabs.New()
	_, _ = events.Array("events")
	for _, event := range getEvents(&req, requester) {
		_ = events.ArrayAppend(event.Serialize(), "events")
	}
	return ReturnJSON(c, events)
}
