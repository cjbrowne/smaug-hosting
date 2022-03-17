package idp

import (
	"bitbucket.org/smaug-hosting/services/database"
	"github.com/Masterminds/squirrel"
)

type ServiceRepository struct{}

const tableName = "services"

func (sr ServiceRepository) FindByClientId(clientId string) (*Service, error) {
	sql, params, err := squirrel.
		Select("name").
		From(tableName).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := database.Connection.QueryRowx(sql, params)

	if row.Err() != nil {
		return nil, row.Err()
	}

	srvc := new(Service)

	err = row.StructScan(srvc)

	return srvc, err
}
