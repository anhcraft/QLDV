package request

type HomepageSettingModel struct {
	ActivitySlideshow        []string `json:"activitySlideshow,omitempty"`
	FeaturedUserLimit        uint8    `json:"featuredUserLimit,omitempty"`
	FeaturedAchievementLimit uint8    `json:"featuredAchievementLimit,omitempty"`
}
