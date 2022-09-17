package models

import "github.com/Jeffail/gabs/v2"

type ContestScoreRank struct {
	Rank  float64
	Count float64
}

func (a *ContestScoreRank) serialize() *gabs.Container {
	res := gabs.New()
	_, _ = res.Set(a.Rank, "rank")
	_, _ = res.Set(a.Count, "count")
	return res
}
