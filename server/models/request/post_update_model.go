package request

type PostUpdateModel struct {
	Title    string `json:"title,omitempty"`
	Headline string `json:"headline,omitempty"`
	Content  string `json:"content,omitempty"`
	Hashtag  string `json:"hashtag,omitempty"`
	Privacy  uint8  `json:"privacy,omitempty"`
}
