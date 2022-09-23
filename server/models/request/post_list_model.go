package request

type PostListModel struct {
	Limit          uint8    `query:"limit,omitempty"`
	FilterHashtags []string `query:"filter-hashtags,omitempty"`
	BelowId        uint32   `query:"below-id,omitempty"`
	SortBy         string   `query:"sort-by,omitempty"`
	LowerThan      uint     `query:"lower-than,omitempty"`
}
