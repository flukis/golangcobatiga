package database

import (
	"pos-services/app/queries"
	"pos-services/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Queries struct {
	*queries.UserQueries
	*queries.AuthQueries
	*queries.UserLocationQueries
}

func OpenDBConnection() (*Queries, error) {
	var (
		db  *gorm.DB
		err error
	)

	urlConn, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	db, err = gorm.Open(postgres.Open(urlConn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:         &queries.UserQueries{DB: db},
		AuthQueries:         &queries.AuthQueries{DB: db},
		UserLocationQueries: &queries.UserLocationQueries{DB: db},
	}, nil
}
