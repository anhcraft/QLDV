package main

import "github.com/Jeffail/gabs/v2"

type ContestSession struct {
	ContestID               string  `gorm:"primaryKey"`
	Contest                 Contest `gorm:"constraint:OnDelete:CASCADE;"`
	UserID                  string  `gorm:"primaryKey"`
	User                    User    `gorm:"constraint:OnDelete:CASCADE;"`
	StartTime               int64
	EndTime                 int64
	QuestionSheet           string
	AnswerSheet             string
	ExpectedAnswerSheet     string
	LastAnswerSubmittedTime int64
	Finished                bool
}

func (a *ContestSession) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.UserID, "userId")
	_, _ = res.Set(a.ContestID, "contestId")
	_, _ = res.Set(a.StartTime, "startTime")
	_, _ = res.Set(a.EndTime, "endTime")
	_, _ = res.Set(a.QuestionSheet, "questionSheet")
	_, _ = res.Set(a.AnswerSheet, "answerSheet")
	_, _ = res.Set(a.LastAnswerSubmittedTime, "lastAnswerSubmittedTime")
	_, _ = res.Set(a.Finished, "finished")
	if a.Finished {
		_, _ = res.Set(a.ExpectedAnswerSheet, "expectedAnswerSheet")
	}
	return res
}
