package storage

import (
	"das/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB

func init() {
	dsn := os.Getenv("sql")
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while connecting to the database")
	}
	err = db_.AutoMigrate(&models.User{}, &models.AnnualRank{}, &models.Achievement{}, &models.Post{}, &models.Attachment{}, &models.Event{}, &models.PostStat{}, &models.Contest{}, &models.ContestSession{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while migrating tables")
	}
	Db = db_
}
