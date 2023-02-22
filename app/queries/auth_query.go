package queries

import (
	"pos-services/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type AuthQueries struct {
	*gorm.DB
}

func (q *AuthQueries) CreateSession(p *utils.PayloadSession) error {
	res := q.DB.Create(p)
	if res.Error != nil {
		if res.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return res.Error
	}
	return nil
}

func (q *AuthQueries) GetSessionById(id string) (*utils.PayloadSession, error) {
	var payload utils.PayloadSession
	res := q.DB.First(&payload, id)
	if res.Error != nil {
		return &payload, res.Error
	}
	return &payload, nil
}

func (q *AuthQueries) GetSessionByEmail(email string) (*utils.PayloadSession, error) {
	var payload utils.PayloadSession
	res := q.DB.Where("email = ?", email).First(&payload)
	if res.Error != nil {
		return &payload, res.Error
	}
	return &payload, nil
}

func (q *AuthQueries) DeleteSession(id string, now time.Time) error {
	res := q.DB.Model(&utils.PayloadSession{}).Where("id = ?", id).Update("expired_at", now)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
