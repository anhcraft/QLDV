package main

import (
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func editOrCreateContest(eventId string, answers bool, questions uint8, time uint32, sheet string, info string) interface{} {
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

func getContestSession(user string, contest string) *ContestSession {
	var contestSession ContestSession
	result := db.Take(&contestSession, "user_id = ? and contest_id = ?", user, contest)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &contestSession
	}
}

func getContestSessions(contest string, limit int, older int64) []ContestSession {
	var contestSessions []ContestSession
	_ = db.Where("contest_id = ? and last_answer_submitted_time < ?", contest, older).Order("last_answer_submitted_time desc").Limit(limit).Find(&contestSessions)
	return contestSessions
}

func createContestSession(user string, contest string, limitTime uint32, questionSheet string, answerSheet string, expectedAnswerSheet string) *ContestSession {
	t := time.Now().UnixMilli()
	c := &ContestSession{
		ContestID:               contest,
		UserID:                  user,
		StartTime:               t,
		EndTime:                 t + int64(limitTime),
		QuestionSheet:           questionSheet,
		AnswerSheet:             answerSheet,
		ExpectedAnswerSheet:     expectedAnswerSheet,
		LastAnswerSubmittedTime: 0,
		Finished:                false,
	}
	_ = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&c)
	return c
}

func submitContestSession(user string, contest string, answerSheet string, saveOnly bool) bool {
	t := time.Now().UnixMilli()
	c := &ContestSession{
		AnswerSheet:             answerSheet,
		LastAnswerSubmittedTime: t,
		Finished:                !saveOnly,
	}
	tx := db.Model(c).Where("user_id = ? and contest_id = ? and finished = ? and start_time <= ? and end_time + 20000 >= ?", user, contest, false, t, t).Select("answer_sheet", "last_answer_submitted_time", "finished").Updates(c)
	return tx.RowsAffected > 0
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
		LimitQuestions   uint8  `json:"limit_questions,omitempty"`
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
	res = contest.serialize(user.Admin)
	return c.SendString(res.String())
}

func contestSessionGetRouteHandler(c *fiber.Ctx) error {
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
	contestSession := getContestSession(email, id)
	if contestSession == nil {
		_, _ = res.Set("ERR_UNKNOWN_CONTEST_SESSION", "error")
		return c.SendString(res.String())
	}
	res = contestSession.serialize()
	return c.SendString(res.String())
}

func contestSessionListRouteHandler(c *fiber.Ctx) error {
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
	limit, err1 := strconv.Atoi(c.Query("limit", ""))
	if err1 != nil || limit > 50 {
		limit = 50
	}
	older, err2 := strconv.ParseInt(c.Query("older", ""), 10, 64)
	if err2 != nil {
		older = time.Now().UnixMilli()
	}
	_, _ = res.Array("contestSessions")
	for _, ev := range getContestSessions(c.Query("contest", ""), limit, older) {
		cont := ev.serialize()
		_ = res.ArrayAppend(cont, "contestSessions")
	}
	return c.SendString(res.String())
}

func contestSessionSubmitRouteHandler(c *fiber.Ctx) error {
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

	payload := struct {
		Id          string `json:"id,omitempty"`
		AnswerSheet string `json:"answer_sheet,omitempty"`
		SaveOnly    bool   `json:"save_only,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	if getContest(payload.Id) == nil {
		_, _ = res.Set("ERR_UNKNOWN_CONTEST", "error")
		return c.SendString(res.String())
	}

	_, _ = res.Set(submitContestSession(email, payload.Id, payload.AnswerSheet, payload.SaveOnly), "success")
	return c.SendString(res.String())
}

func contestSessionJoinRouteHandler(c *fiber.Ctx) error {
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

	payload := struct {
		Id string `questionSheetJSON:"id,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	ev := getEvent(payload.Id)
	t := time.Now().UnixMilli()
	if ev.StartDate > t || t > ev.EndDate {
		_, _ = res.Set("ERR_EVENT_UNAVAILABLE", "error")
		return c.SendString(res.String())
	}
	if (ev.Privacy&2) == 2 && !(user.Mod || user.Admin) {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}
	if (ev.Privacy&4) == 4 && !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	contest := getContest(payload.Id)
	if contest == nil {
		_, _ = res.Set("ERR_UNKNOWN_CONTEST", "error")
		return c.SendString(res.String())
	}
	if !contest.AcceptingAnswers {
		_, _ = res.Set("ERR_CONTEST_CLOSED", "error")
		return c.SendString(res.String())
	}

	if getContestSession(email, payload.Id) != nil {
		_, _ = res.Set("ERR_CONTEST_ATTENDED", "error")
		return c.SendString(res.String())
	}

	dataSheet, _ := gabs.ParseJSON([]byte(contest.DataSheet))
	in := dataSheet.Children()
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})
	expectedAnswerSheet, _ := gabs.ParseJSON([]byte("[]"))
	questionSheet, _ := gabs.ParseJSON([]byte("[]"))
	for i := 0; i < int(contest.LimitQuestions); i++ {
		_ = expectedAnswerSheet.ArrayAppend(in[i].Path("answer").Data().(float64))
		_ = in[i].Delete("answer")
		_ = questionSheet.ArrayAppend(in[i])
	}
	questionSheetJSON, err := questionSheet.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	expectedAnswerSheetJSON, err := expectedAnswerSheet.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	answerSheet, _ := gabs.ParseJSON([]byte("[]"))
	for i := 0; i < int(contest.LimitQuestions); i++ {
		_ = answerSheet.ArrayAppend(-1)
	}
	answerSheetJSON, err := answerSheet.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	contestSession := createContestSession(email, payload.Id, contest.LimitTime, string(questionSheetJSON), string(answerSheetJSON), string(expectedAnswerSheetJSON))
	res = contestSession.serialize()
	return c.SendString(res.String())
}
