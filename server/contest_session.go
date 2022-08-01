package main

import (
	"github.com/Jeffail/gabs/v2"
	"time"
)

type ContestSession struct {
	ID                      string `gorm:"primaryKey"`
	ContestID               int
	UserID                  string
	StartTime               int64
	EndTime                 int64
	QuestionSheet           string
	AnswerSheet             string
	ExpectedAnswerSheet     string
	LastAnswerSubmittedTime int64
	Finished                bool
	Score                   float32
	AnswerAccuracy          string
}

func (a *ContestSession) serialize(showFull bool) *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.ID, "id")
	_, _ = res.Set(a.UserID, "userId")
	_, _ = res.Set(a.ContestID, "contestId")
	if showFull {
		_, _ = res.Set(a.StartTime, "startTime")
		_, _ = res.Set(a.EndTime, "endTime")
		_, _ = res.Set(a.QuestionSheet, "questionSheet")
		_, _ = res.Set(a.AnswerSheet, "answerSheet")
		_, _ = res.Set(a.LastAnswerSubmittedTime, "lastAnswerSubmittedTime")
	}
	_, _ = res.Set(a.Finished, "finished")
	if a.Finished {
		if showFull {
			//_, _ = res.Set(a.ExpectedAnswerSheet, "expectedAnswerSheet")
			_, _ = res.Set(a.AnswerAccuracy, "answerAccuracy")
		}
		_, _ = res.Set(a.Score, "score")
	} else {
		_, _ = res.Set(time.Now().UnixMilli() >= a.EndTime, "finished")
	}
	return res
}
