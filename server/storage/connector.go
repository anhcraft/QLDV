package storage

import (
	"das/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var Db *gorm.DB

func init() {
	dsn := os.Getenv("sql")
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while connecting to the database")
	}
	err = db_.AutoMigrate(&models.User{}, &models.AnnualRank{}, &models.Achievement{}, &models.Post{}, &models.PostStat{}, &models.Attachment{}, &models.Event{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while migrating tables")
	}
	Db = db_
}
