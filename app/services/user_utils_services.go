package services

import (
	"net/http"
	"pos-services/app/models"
	"pos-services/platform/database"
)

func AddAddress(id string, a *models.UserLocation) (models.UserLocation, ServiceReturn) {
	// add address to database address
	// check if user is already registered
	db, err := database.OpenDBConnection()
	if err != nil {
		return models.UserLocation{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	a.UserID = id

	user, err := db.GetById(id)
	if err != nil {
		return models.UserLocation{}, ServiceReturn{
			HttpStatusCode: http.StatusNotFound,
			Err:            err,
		}
	}

	err = db.AddAddress(user, a)
	if err != nil {
		return models.UserLocation{}, ServiceReturn{
			HttpStatusCode: http.StatusNotFound,
			Err:            err,
		}
	}

	return *a, ServiceReturn{
		HttpStatusCode: http.StatusCreated,
		Err:            nil,
	}
}
