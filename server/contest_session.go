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
		//_, _ = res.Set(a.ExpectedAnswerSheet, "expectedAnswerSheet")
		answerSheet, _ := gabs.ParseJSON([]byte(a.AnswerSheet))
		answerSheetList := answerSheet.Children()
		expectedAnswerSheet, _ := gabs.ParseJSON([]byte(a.ExpectedAnswerSheet))
		expectedAnswerSheetList := expectedAnswerSheet.Children()
		_, _ = res.Array("answerAccuracy")
		for i := 0; i < len(answerSheetList); i++ {
			_ = res.ArrayAppend(answerSheetList[i].Data().(float64) == expectedAnswerSheetList[i].Data().(float64), "answerAccuracy")
		}
	}
	return res
}
