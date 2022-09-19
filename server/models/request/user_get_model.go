package request

type UserGetModel struct {
	Profile      bool `query:"profile"`
	Achievements bool `query:"achievements"`
	AnnualRanks  bool `query:"annual-ranks"`
}
