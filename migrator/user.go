package main

type User struct {
	Class     string `json:"class"`
	Name      string `json:"name"`
	Birth     int64  `json:"birth"`
	Gender    bool   `json:"gender"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Certified bool   `json:"certified"`
}
