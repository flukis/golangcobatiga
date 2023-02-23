package migrations

import (
	"pos-services/app/models"
	"pos-services/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() error {
	urlConn, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(urlConn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(
		&models.Users{},
		&models.UserLocation{},
		&utils.PayloadSession{},
	)
	return nil
}
