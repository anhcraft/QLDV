package models

type Settings struct {
	Key     string `gorm:"primaryKey"`
	Value   string
	Privacy bool
}
