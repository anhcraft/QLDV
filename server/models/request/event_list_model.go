package request

type EventListModel struct {
	Limit     uint8  `query:"limit,omitempty"`
	BelowId   uint32 `query:"below-id,omitempty"`
	BeginDate uint64 `query:"begin-date,omitempty"`
	EndDate   uint64 `query:"end-date,omitempty"`
}
