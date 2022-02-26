package main

import (
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func editOrCreateContest(eventId string, answers bool, questions uint16, time uint32, sheet string, info string) interface{} {
	contest := &Contest{
		AcceptingAnswers: answers,
		LimitQuestions:   questions,
		LimitTime:        time,
		DataSheet:        sheet,
		EventID:          eventId,
		Info:             info,
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "event_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"accepting_answers", "limit_questions", "limit_time", "data_sheet", "info"}),
	}).Create(&contest)
	return &contest
}

func removeContest(id string) bool {
	var contest Contest
	db.Where("event_id = ?", id).Delete(&contest)
	return true
}

func getContest(id string) *Contest {
	var contest Contest
	result := db.Take(&contest, "event_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &contest
	}
}

func contestChangeRouteHandler(c *fiber.Ctx) error {
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

	payload := struct {
		Id               string `json:"id,omitempty"`
		AcceptingAnswers bool   `json:"accepting_answers,omitempty"`
		LimitQuestions   uint16 `json:"limit_questions,omitempty"`
		LimitTime        uint32 `json:"limit_time,omitempty"`
		DataSheet        string `json:"data_sheet,omitempty"`
		Info             string `json:"info,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	if getEvent(payload.Id) == nil {
		_, _ = res.Set("ERR_UNKNOWN_EVENT", "error")
		return c.SendString(res.String())
	}

	payload.Info = ugcPolicy.Sanitize(payload.Info)
	_ = editOrCreateContest(payload.Id, payload.AcceptingAnswers, payload.LimitQuestions, payload.LimitTime, payload.DataSheet, payload.Info)
	_, _ = res.Set(true, "success")
	return c.SendString(res.String())
}

func contestRemoveRouteHandler(c *fiber.Ctx) error {
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
	_, _ = res.Set(removeContest(id), "success")
	return c.SendString(res.String())
}

func contestGetRouteHandler(c *fiber.Ctx) error {
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

	id := c.Query("id", "")
	if id == "" {
		_, _ = res.Set("ERR_INVALID_CONTEST_ID", "error")
		return c.SendString(res.String())
	}
	contest := getContest(id)
	if contest == nil {
		_, _ = res.Set("ERR_UNKNOWN_CONTEST", "error")
		return c.SendString(res.String())
	}
	res = contest.serialize()
	return c.SendString(res.String())
}
