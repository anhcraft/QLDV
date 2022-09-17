package request

type PostListModel struct {
	Limit          uint8    `json:"limit,omitempty"`
	FilterHashtags []string `json:"filterHashtags,omitempty"`
	BelowId        uint32   `json:"belowId,omitempty"`
	SortBy         string   `json:"sortBy,omitempty"`
	LowerThan      uint     `json:"lowerThan,omitempty"`
}
