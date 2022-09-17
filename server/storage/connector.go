package storage

import (
	"das/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB

func init() {
	dsn := os.Getenv("sql")
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting database: %v\n", err)
	}
	err = db_.AutoMigrate(&models.User{}, &models.AnnualRank{}, &models.Achievement{}, &models.Post{}, &models.Attachment{}, &models.Event{}, &models.PostStat{}, &models.Contest{}, &models.ContestSession{})
	if err != nil {
		log.Fatalf("error migrating: %v\n", err)
	}
	Db = db_
}
