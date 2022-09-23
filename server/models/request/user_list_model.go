package request

type UserListModel struct {
	Limit       uint8  `query:"limit,omitempty"`
	FilterName  string `query:"filter-name,omitempty"`
	FilterClass string `query:"filter-class,omitempty"`
	FilterEmail string `query:"filter-email,omitempty"`
	FilterRole  uint8  `query:"filter-role,omitempty"`
	BelowId     uint16 `query:"below-id,omitempty"`
}
