package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"math/rand"
	"strings"
	"time"
)

func editOrCreateContest(eventId string, answers bool, limitQuestions uint8, limitTime uint32, limitSessions uint8, sheet string, info string) interface{} {
	contest := &Contest{
		AcceptingAnswers: answers,
		LimitQuestions:   limitQuestions,
		LimitTime:        limitTime,
		LimitSessions:    limitSessions,
		DataSheet:        sheet,
		EventID:          eventId,
		Info:             info,
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "event_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"accepting_answers", "limit_questions", "limit_time", "limit_sessions", "data_sheet", "info"}),
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

func containString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getContestSessions(contest string, limit int, offset int, attendant string, requireFinished bool, order []string) []ContestSession {
	var contestSessions []ContestSession
	a := db.Where("contest_id = ?", contest)
	if len(attendant) > 0 {
		a = a.Where("LOWER(`user_id`) like ?", "%"+attendant+"%")
	}
	if requireFinished {
		t := time.Now().UnixMilli()
		a = a.Where("(? - `end_time` > 0) or (finished = ?)", t, true)
	}
	if containString(order, "final-score") {
		a = a.Order("score desc")
	}
	if containString(order, "session-time") {
		a = a.Order("(`last_answer_submitted_time` - `start_time`) desc")
	} else { // session-time can not be combined with sorting by latest submitted time
		a = a.Order("last_answer_submitted_time desc")
	}
	a = a.Offset(offset).Limit(limit).Find(&contestSessions)
	return contestSessions
}

func createContestSession(user string, contest string, limitTime uint32, questionSheet string, answerSheet string, expectedAnswerSheet string) *ContestSession {
	t := time.Now().UnixMilli()
	hash := sha256.New()
	hash.Write([]byte(user + contest + time.Now().String()))
	id := hex.EncodeToString(hash.Sum(nil))
	c := &ContestSession{
		ID:                      id,
		ContestID:               contest,
		UserID:                  user,
		StartTime:               t,
		EndTime:                 t + int64(limitTime),
		QuestionSheet:           questionSheet,
		AnswerSheet:             answerSheet,
		ExpectedAnswerSheet:     expectedAnswerSheet,
		LastAnswerSubmittedTime: 0,
		Finished:                false,
		Score:                   0,
		AnswerAccuracy:          "[]",
	}
	_ = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&c)
	return c
}

func submitContestSession(id string, answerSheet string, saveOnly bool) bool {
	t := time.Now().UnixMilli()
	score := float32(0)
	answerAccuracy := "[]"
	if !saveOnly {
		var contestSession ContestSession
		result := db.Take(&contestSession, "id = ? and finished = ? and start_time <= ? and end_time + 20000 >= ?", id, false, t, t)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		answerSheetJSON, err := gabs.ParseJSON([]byte(answerSheet))
		if err != nil {
			return false
		}
		answerSheetList := answerSheetJSON.Children()
		expectedAnswerSheetJSON, err := gabs.ParseJSON([]byte(contestSession.ExpectedAnswerSheet))
		if err != nil {
			return false
		}
		expectedAnswerSheetList := expectedAnswerSheetJSON.Children()
		answerAccuracyJSON, err := gabs.ParseJSON([]byte("[]"))
		if err != nil {
			return false
		}
		j := float32(0)
		for i := 0; i < len(answerSheetList); i++ {
			b := answerSheetList[i].Data().(float64) == expectedAnswerSheetList[i].Data().(float64)
			_ = answerAccuracyJSON.ArrayAppend(b)
			if b {
				j++
			}
		}
		score = j / float32(len(expectedAnswerSheetList)) * 10.0
		answerAccuracyBytes, err := answerAccuracyJSON.MarshalJSON()
		if err != nil {
			return false
		}
		answerAccuracy = string(answerAccuracyBytes)
	}
	c := &ContestSession{
		AnswerSheet:             answerSheet,
		LastAnswerSubmittedTime: t,
		Finished:                !saveOnly,
		Score:                   score,
		AnswerAccuracy:          answerAccuracy,
	}
	tx := db.Model(c).Where("id = ? and finished = ? and start_time <= ? and end_time + 20000 >= ?", id, false, t, t)
	tx = tx.Select("answer_sheet", "last_answer_submitted_time", "finished", "score", "answer_accuracy").Updates(c)
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
		LimitSessions    uint8  `json:"limit_sessions,omitempty"`
		DataSheet        string `json:"data_sheet,omitempty"`
		Info             string `json:"info,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}
	payload.Info = strings.TrimSpace(payload.Info)
	if payload.LimitSessions > 30 {
		payload.LimitSessions = 30
	}

	if getEvent(payload.Id) == nil {
		_, _ = res.Set("ERR_UNKNOWN_EVENT", "error")
		return c.SendString(res.String())
	}

	payload.Info = ugcPolicy.Sanitize(payload.Info)
	_ = editOrCreateContest(payload.Id, payload.AcceptingAnswers, payload.LimitQuestions, payload.LimitTime, payload.LimitSessions, payload.DataSheet, payload.Info)
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

	payload := struct {
		Contest         string   `json:"contest,omitempty"`
		Limit           int      `json:"limit,omitempty"`
		Offset          int      `json:"offset,omitempty"`
		FilterAttendant string   `json:"filter_attendant,omitempty"`
		FilterFinished  bool     `json:"filter_finished,omitempty"`
		SortBy          []string `json:"sort_by,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}
	payload.FilterAttendant = strings.TrimSpace(payload.FilterAttendant)
	if payload.Limit > 50 {
		payload.Limit = 50
	}
	if payload.Offset < 0 {
		payload.Offset = 0
	}
	if !user.Admin {
		payload.FilterAttendant = user.Email
	}
	_, _ = res.Array("contestSessions")
	for _, ev := range getContestSessions(payload.Contest, payload.Limit, payload.Offset, payload.FilterAttendant, payload.FilterFinished, payload.SortBy) {
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

	_, _ = res.Set(submitContestSession(payload.Id, payload.AnswerSheet, payload.SaveOnly), "success")
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
		Id string `json:"id,omitempty"`
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

	joinedSessions := getContestSessions(payload.Id, int(contest.LimitSessions), 0, email, false, []string{})
	if len(joinedSessions) >= int(contest.LimitSessions) {
		_, _ = res.Set("ERR_CONTEST_ATTENDED_MAX", "error")
		return c.SendString(res.String())
	}
	for i := 0; i < len(joinedSessions); i++ {
		if !joinedSessions[i].Finished && joinedSessions[i].EndTime >= t {
			_, _ = res.Set("ERR_CONTEST_ATTENDED", "error")
			return c.SendString(res.String())
		}
	}

	dataSheet, _ := gabs.ParseJSON([]byte(contest.DataSheet))
	in := dataSheet.Children()
	if len(in) < int(contest.LimitQuestions) {
		_, _ = res.Set("ERR_CONTEST_DATA_INSUFFICIENT", "error")
		return c.SendString(res.String())
	}
	rand.Seed(t)
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
	_, _ = res.Set(contestSession.ID, "id")
	return c.SendString(res.String())
}
