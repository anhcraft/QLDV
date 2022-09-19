package request

type EventListModel struct {
	Limit     uint8  `json:"limit,omitempty"`
	BelowId   uint32 `json:"belowId,omitempty"`
	BeginDate uint64 `json:"beginDate,omitempty"`
	EndDate   uint64 `json:"endDate,omitempty"`
}
