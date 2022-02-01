package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

func getEvents(limit int, older int64, fromDate int64, toDate int64) []Event {
	var events []Event
	a := db.Where("date < ?", older)
	if fromDate != 0 {
		a = a.Where("start_date >= ? or end_date >= ?", fromDate, fromDate)
	}
	if toDate != 0 {
		a = a.Where("start_date <= ? or end_date <= ?", toDate, toDate)
	}
	a.Order("date desc").Limit(limit).Find(&events)
	return events
}

func removeEvent(id string) bool {
	var event Event
	db.Where("id = ?", id).Delete(&event)
	return true
}

func getEvent(id string) *Event {
	var event Event
	result := db.Take(&event, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &event
	}
}

func editOrCreateEvent(id string, title string, startDate int64, endDate int64, privacy uint8) *Event {
	if id == "" {
		hash := sha256.New()
		hash.Write([]byte(id + title + time.Now().String()))
		md := hash.Sum(nil)
		id = hex.EncodeToString(md)
	}
	event := Event{
		ID:        id,
		Title:     title,
		StartDate: startDate,
		EndDate:   endDate,
		Date:      time.Now().UnixMilli(),
		Privacy:   privacy,
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "start_date", "end_date", "privacy"}),
	}).Create(&event)
	return &event
}

func eventGetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	id := c.Query("id", "")
	if id == "" {
		_, _ = res.Set("ERR_INVALID_POST_ID", "error")
		return c.SendString(res.String())
	}
	event := getEvent(id)
	if event == nil {
		_, _ = res.Set("ERR_UNKNOWN_POST", "error")
		return c.SendString(res.String())
	}

	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	var user *User = nil
	if success {
		user = getProfile(email)
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
	return c.SendString(res.String())
}

func eventListRouteHandler(c *fiber.Ctx) error {
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	var user *User = nil
	if success {
		user = getProfile(email)
	}

	res := gabs.New()
	limit, err1 := strconv.Atoi(c.Query("limit", ""))
	if err1 != nil || limit > 50 {
		limit = 50
	}
	older, err2 := strconv.ParseInt(c.Query("older", ""), 10, 64)
	if err2 != nil {
		older = time.Now().UnixMilli()
	}
	fromDate, err3 := strconv.ParseInt(c.Query("from-date", ""), 10, 64)
	if err3 != nil {
		fromDate = 0
	}
	toDate, err4 := strconv.ParseInt(c.Query("to-date", ""), 10, 64)
	if err4 != nil {
		toDate = 0
	}
	_, _ = res.Array("events")
	for _, ev := range getEvents(limit, older, fromDate, toDate) {
		if (ev.Privacy&1) == 1 && user == nil {
			continue
		}
		if (ev.Privacy&2) == 2 && (user == nil || !(user.Mod || user.Admin)) {
			continue
		}
		if (ev.Privacy&4) == 4 && (user == nil || !user.Admin) {
			continue
		}
		_ = res.ArrayAppend(ev.serialize(), "events")
	}
	return c.SendString(res.String())
}

func eventRemoveRouteHandler(c *fiber.Ctx) error {
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
	_, _ = res.Set(removeEvent(id), "success")
	return c.SendString(res.String())
}

func eventChangeRouteHandler(c *fiber.Ctx) error {
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
		Id        string `json:"id,omitempty"`
		Title     string `json:"title,omitempty"`
		StartDate int64  `json:"start_date,omitempty"`
		EndDate   int64  `json:"end_date,omitempty"`
		Privacy   uint8  `json:"privacy,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	if len(payload.Title) < 5 {
		_, _ = res.Set("ERR_POST_TITLE_MIN", "error")
		return c.SendString(res.String())
	} else if len(payload.Title) > 300 {
		_, _ = res.Set("ERR_POST_TITLE_MAX", "error")
		return c.SendString(res.String())
	}

	if payload.StartDate > payload.EndDate {
		_, _ = res.Set("ERR_DATE_RANGE", "error")
		return c.SendString(res.String())
	}

	p := editOrCreateEvent(payload.Id, payload.Title, payload.StartDate, payload.EndDate, payload.Privacy)
	_, _ = res.Set(true, "success")
	_, _ = res.Set(p.ID, "id")
	return c.SendString(res.String())
}
