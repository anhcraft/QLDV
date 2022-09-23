package request

type EventUpdateModel struct {
	Title     string `json:"title,omitempty"`
	BeginDate uint64 `json:"beginDate,omitempty"`
	EndDate   uint64 `json:"endDate,omitempty"`
	Privacy   uint8  `json:"privacy,omitempty"`
}
