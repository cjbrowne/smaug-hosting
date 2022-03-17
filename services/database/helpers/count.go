package helpers

import (
	"bitbucket.org/smaug-hosting/services/database"
	"errors"
	"github.com/Masterminds/squirrel"
)

var (
	ErrNoRowsReturned = errors.New("no rows returned from sql query")
)

func Count(table string, where string, args ...interface{}) (uint64, error) {
	builder := squirrel.Select("COUNT(*)").From(table)

	if where != "" {
		builder = builder.Where(where, args...)
	}

	sql, params, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := database.Connection.Query(sql, params...)
	if err != nil {
		return 0, err
	}

	var count uint64

	if rows.Next() {
		err = rows.Scan(&count)
	} else {
		err = ErrNoRowsReturned
	}

	return count, err
}
