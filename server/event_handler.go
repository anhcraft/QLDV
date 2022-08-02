package main

import (
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

func getEvents(limit int, belowId int, beginDate int64, endDate int64) []Event {
	var events []Event
	a := db.Limit(limit).Order("end_date desc, begin_date desc, id desc")
	if beginDate > 0 && endDate > 0 && beginDate <= endDate {
		a = a.Where("begin_date <= ? and end_date >= ?", endDate, beginDate)
	} else if beginDate > 0 && endDate == 0 {
		a = a.Where("begin_date > ?", beginDate)
	} else if beginDate == 0 && endDate > 0 {
		a = a.Where("end_date < ?", endDate)
	}
	if belowId > 0 {
		a = a.Where("id < ?", belowId)
	}
	a.Find(&events)
	return events
}

func removeEvent(id int) bool {
	var event Event
	db.Where("id = ?", id).Delete(&event)
	return true
}

func getEvent(id int) *Event {
	var event Event
	result := db.Take(&event, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &event
	}
}

func editOrCreateEvent(id int, title string, beginDate int64, endDate int64, privacy uint8) *Event {
	event := Event{
		Title:     title,
		Link:      GenerateLinkFromTitle(title),
		BeginDate: beginDate,
		EndDate:   endDate,
		Date:      time.Now().UnixMilli(),
		Privacy:   privacy,
	}
	if id > 0 {
		event.ID = id
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "link", "begin_date", "end_date", "privacy"}),
	}).Create(&event)
	return &event
}

func eventGetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		_, _ = res.Set("ERR_INVALID_EVENT_ID", "error")
		return c.SendString(res.String())
	}
	event := getEvent(id)
	if event == nil {
		_, _ = res.Set("ERR_UNKNOWN_EVENT", "error")
		return c.SendString(res.String())
	}

	token := c.Get("token")
	success, emailOrError := getEmailFromToken(token, c.UserContext())
	var user *User = nil
	if success {
		user = getProfile(emailOrError)
	}
	if (event.Privacy&1) == 1 && user == nil {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}
	if (event.Privacy&2) == 2 && (user == nil || !(user.Mod || user.Admin)) {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}
	if (event.Privacy&4) == 4 && (user == nil || !user.Admin) {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	res = event.serialize()
	contest := getContest(id)
	if contest != nil && ((user != nil && user.Admin) || contest.AcceptingAnswers) {
		_, _ = res.Set(contest.serialize(user != nil && user.Admin), "contest")
	}
	return c.SendString(res.String())
}

func eventListRouteHandler(c *fiber.Ctx) error {
	token := c.Get("token")
	success, emailOrError := getEmailFromToken(token, c.UserContext())
	var user *User = nil
	if success {
		user = getProfile(emailOrError)
	}

	res := gabs.New()
	limit, err1 := strconv.Atoi(c.Query("limit", ""))
	if err1 != nil || limit > 50 {
		limit = 50
	} else if limit < 1 {
		limit = 1
	}
	belowId, err2 := strconv.Atoi(c.Query("below-id", ""))
	if err2 != nil {
		belowId = 0
	}
	beginDate, err3 := strconv.ParseInt(c.Query("begin-date", ""), 10, 64)
	if err3 != nil {
		beginDate = 0
	}
	endDate, err4 := strconv.ParseInt(c.Query("end-date", ""), 10, 64)
	if err4 != nil {
		endDate = 0
	}
	_, _ = res.Array("events")
	for _, ev := range getEvents(limit, belowId, beginDate, endDate) {
		if (ev.Privacy&1) == 1 && user == nil {
			continue
		}
		if (ev.Privacy&2) == 2 && (user == nil || !(user.Mod || user.Admin)) {
			continue
		}
		if (ev.Privacy&4) == 4 && (user == nil || !user.Admin) {
			continue
		}
		cont := ev.serialize()
		contest := getContest(ev.ID)
		if contest != nil && ((user != nil && user.Admin) || contest.AcceptingAnswers) {
			_, _ = cont.Set(contest.serialize(user != nil && user.Admin), "contest")
		}
		_ = res.ArrayAppend(cont, "events")
	}
	return c.SendString(res.String())
}

func eventRemoveRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(emailOrError, "error")
		return c.SendString(res.String())
	}
	user := getProfile(emailOrError)
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
		_, _ = res.Set("ERR_INVALID_EVENT_ID", "error")
		return c.SendString(res.String())
	}
	_, _ = res.Set(removeEvent(id), "success")
	return c.SendString(res.String())
}

func eventChangeRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(emailOrError, "error")
		return c.SendString(res.String())
	}
	user := getProfile(emailOrError)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	payload := struct {
		Id        int    `json:"id,omitempty"`
		Title     string `json:"title,omitempty"`
		BeginDate int64  `json:"begin_date,omitempty"`
		EndDate   int64  `json:"end_date,omitempty"`
		Privacy   uint8  `json:"privacy,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Title = ugcPolicy.Sanitize(payload.Title)

	if len(payload.Title) < 5 {
		_, _ = res.Set("ERR_EVENT_TITLE_MIN", "error")
		return c.SendString(res.String())
	} else if len(payload.Title) > 300 {
		_, _ = res.Set("ERR_EVENT_TITLE_MAX", "error")
		return c.SendString(res.String())
	}

	if payload.BeginDate > payload.EndDate {
		_, _ = res.Set("ERR_DATE_RANGE", "error")
		return c.SendString(res.String())
	}

	p := editOrCreateEvent(payload.Id, payload.Title, payload.BeginDate, payload.EndDate, payload.Privacy)
	_, _ = res.Set(true, "success")
	_, _ = res.Set(p.ID, "id")
	return c.SendString(res.String())
}
