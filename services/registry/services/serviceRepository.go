package services

import (
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/database/helpers"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ServiceRepository struct {
}

type ServiceSearchQuery struct {
	Page       uint64
	PageSize   uint64
	Capability string
}

type ServiceSearchResult struct {
	Total    uint64
	Page     uint64
	PageSize uint64
	Services []Service
	Error    error `json:"-"`
}

const tableName = "services"

var (
	ErrNotFound = errors.New("service not found in database")
)

func (r ServiceRepository) All() []Service {
	cols := helpers.GetColsForType(Service{})

	sql, params, err := squirrel.Select(cols...).From(tableName).ToSql()
	if err != nil {
		logrus.Errorf("Could not build database query: %s", err)
		return nil
	}

	rows, err := database.Connection.Query(sql, params...)
	if err != nil {
		logrus.Errorf("Could not run database query: %s", err)
		return nil
	}

	services := make([]Service, 0)
	for rows.Next() {
		service := Service{}

		err = sqlx.StructScan(rows, &service)
		if err != nil {
			logrus.Errorf("Could not parse row into struct: %s", err)
			continue
		}

		services = append(services, service)
	}

	return services
}

func (r ServiceRepository) Find(query ServiceSearchQuery) ServiceSearchResult {
	cols := helpers.GetColsForType(Service{})

	var total uint64 = 0
	var err error = nil

	if query.Capability != "" {
		total, err = helpers.Count(tableName, "? IN capability", query.Capability)
	} else {
		total, err = helpers.Count(tableName, "", nil)
	}

	if err != nil {
		logrus.Errorf("Could not get total services with capability %s: %s", query.Capability, err)
		return ServiceSearchResult{Error: err}
	}

	queryBuilder := squirrel.
		Select(cols...).
		From("services").
		Offset(query.Page).
		Limit(query.PageSize)

	if query.Capability != "" {
		queryBuilder = queryBuilder.Where("? IN capability", query.Capability)
	}

	sql, args, err := queryBuilder.ToSql()

	logrus.Tracef("Built SQL: %s", sql)

	if err != nil {
		logrus.Errorf("Could not create SQL query: %s", err)
		return ServiceSearchResult{Error: err}
	}

	rows, err := database.Connection.Query(sql, args...)
	if err != nil {
		logrus.Errorf("Could not execute SQL query: %s", err)
		return ServiceSearchResult{Error: err}
	}

	services := make([]Service, 0)

	err = sqlx.StructScan(rows, &services)
	if err != nil {
		logrus.Errorf("Could not map SQL rows into structs: %s", err)
		return ServiceSearchResult{Error: err}
	}

	return ServiceSearchResult{
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
		Services: services,
		Error:    nil,
	}
}

func (r ServiceRepository) FindOneById(id int) (Service, error) {
	var s Service
	cols := helpers.GetColsForType(Service{})

	sql, params, err := squirrel.Select(cols...).From(tableName).Where("id = ?", id).ToSql()
	if err != nil {
		logrus.Errorf("Could not build SQL query: %s", err)
		return s, err
	}

	rows, err := database.Connection.Query(sql, params)
	if err != nil {
		logrus.Errorf("Could not build SQL query: %s", err)
		return s, err
	}

	if !rows.Next() {
		return s, ErrNotFound
	}

	err = sqlx.StructScan(rows, &s)
	if err != nil {
		logrus.Errorf("Could not scan service into struct: %s", err)
		return s, err
	}

	return s, nil
}

func (r ServiceRepository) Save(s Service) (int64, error) {
	var err error
	var sql string
	var params []interface{}

	if s.Id == 0 {
		sql, params, err = squirrel.Insert(tableName).SetMap(helpers.GetMapForType(s)).ToSql()
	} else {
		sql, params, err = squirrel.Update(tableName).Where("id=?", s.Id).SetMap(helpers.GetMapForType(s)).ToSql()
	}

	if err != nil {
		logrus.Errorf("Could not build SQL query: %s", err)
		return 0, err
	}

	result, err := database.Connection.Exec(sql, params)
	if err != nil {
		logrus.Errorf("Could not save service: %s", err)
		return 0, err
	}

	return result.LastInsertId()
}
