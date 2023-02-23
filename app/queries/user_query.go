package queries

import (
	"pos-services/app/models"

	"gorm.io/gorm"
)

type UserQueries struct {
	*gorm.DB
}

func (q *UserQueries) CreateUser(u *models.Users) error {
	res := q.DB.Create(u)
	if res.Error != nil {
		if res.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return res.Error
	}
	return nil
}

func (q *UserQueries) GetByEmail(email string) (*models.Users, error) {
	var user models.Users
	res := q.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return &user, res.Error
	}
	return &user, nil
}

func (q *UserQueries) GetById(id string) (*models.Users, error) {
	var user models.Users
	res := q.DB.Where("id = ?", id).First(&user)
	if res.Error != nil {
		return &user, res.Error
	}
	return &user, nil
}

func (q *UserQueries) UpdateUser(u *models.Users) error {
	res := q.DB.Save(u)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (q *UserQueries) DeleteUser(id string) error {
	res := q.DB.Delete(&models.Users{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (q *UserQueries) AddAddress(u *models.Users, a *models.UserLocation) error {
	err := q.DB.Model(u).Association("UserLocations").Append([]*models.UserLocation{a})
	if err != nil {
		return err
	}
	return nil
}

func (q *UserQueries) RemoveAddress(u *models.Users, a *models.UserLocation) error {
	err := q.DB.Model(u).Association("UserLocations").Delete([]*models.UserLocation{a})
	if err != nil {
		return err
	}
	return nil
}
