package queries

import (
	"pos-services/app/models"

	"gorm.io/gorm"
)

type UserLocationQueries struct {
	*gorm.DB
}

func (q *UserLocationQueries) GetLocationByUserId(userId string) (*models.UserLocation, error) {
	var user models.UserLocation
	res := q.DB.Where("user_id = ?", userId).First(&user)
	if res.Error != nil {
		return &user, res.Error
	}
	return &user, nil
}

func (q *UserLocationQueries) GetLocationById(id string) (*models.UserLocation, error) {
	var user models.UserLocation
	res := q.DB.First(&user, id)
	if res.Error != nil {
		return &user, res.Error
	}
	return &user, nil
}

func (q *UserLocationQueries) UpdateUserLocation(u *models.UserLocation) error {
	res := q.DB.Save(u)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (q *UserLocationQueries) DeleteUserLocation(id string) error {
	res := q.DB.Delete(&models.UserLocation{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
