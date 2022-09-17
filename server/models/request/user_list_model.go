package request

type UserListModel struct {
	Limit       uint8  `json:"limit,omitempty"`
	FilterName  string `json:"filterName,omitempty"`
	FilterClass string `json:"filterClass,omitempty"`
	FilterEmail string `json:"filterEmail,omitempty"`
	FilterRole  uint8  `json:"filterRole,omitempty"`
	BelowId     uint16 `json:"belowId,omitempty"`
}
