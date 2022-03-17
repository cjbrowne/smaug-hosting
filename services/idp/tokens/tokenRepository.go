package tokens

import (
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/database/helpers"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type TokenRepository struct {
	ctx context.Context
}

const tableName = "tokens"

func (t TokenRepository) Save(token Token) error {
	qry := squirrel.
		Insert(tableName).
		SetMap(token.AsMap()).
		RunWith(database.Connection).
		// upsert
		Suffix("ON DUPLICATE KEY UPDATE token = ?", token.Token)
	res, err := qry.
		Exec()

	if err != nil {
		sql, params, err2 := qry.ToSql()
		logrus.Debugf("SQL: %s [%s] (%s)", sql, params, err2)
		logrus.Errorf("Could not save token: %s", err)
		return err
	}

	if rows, err := res.RowsAffected(); rows < 1 || err != nil {
		logrus.Errorf("Could not save token: %s (rows affected: %d)", err, rows)
		return err
	} else {
		return nil
	}
}

func (t TokenRepository) FindByRefreshToken(refreshToken []byte) (Token, error) {
	var tok Token

	cols := helpers.GetColsForType(Token{})
	qry := squirrel.Select(cols...).From(tableName).Where("refresh = ?", refreshToken)
	sql, params, err := qry.ToSql()

	if err != nil {
		return tok, err
	}

	row := database.Connection.QueryRow(sql, params)
	err = row.Scan(&tok)

	return tok, err
}
