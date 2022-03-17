package users

import (
	"bitbucket.org/smaug-hosting/services/database"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	ctx context.Context
}

const unverifiedTableName = "users"
const verifiedTableName = "verified_users"

var ErrUserNotFound = errors.New("user not found")

func (u UserRepository) FindByEmailAndVerified(email string, verified bool) (*User, error) {
	var usr *User
	cols := getCols(verified)

	tableName := verifiedTableName
	if !verified {
		tableName = unverifiedTableName
	}

	sql, params, err := squirrel.
		Select(cols...).
		From(tableName).
		Where("email = ?", email).
		ToSql()
	if err != nil {
		logrus.Errorf("Could not build SQL query for fetching user: %s", err)
		return usr, err
	}

	rows, err := database.Connection.Queryx(sql, params...)
	if err != nil {
		logrus.Errorf("Could not run SQL query for fetching user: %s", err)
		return usr, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Errorf("Could not close rows: %s", err)
		}
	}()

	if rows.Next() {
		usr = new(User)
		err = rows.StructScan(usr)
		if err != nil {
			logrus.Errorf("Could not scan user into struct: %s", err)
			return usr, err
		} else {
			return usr, nil
		}
	} else {
		// no user found
		return nil, nil
	}
}

func (u UserRepository) FindByEmail(email string) (*User, error) {
	usr, err := u.FindByEmailAndVerified(email, false)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		return usr, nil
	}

	return u.FindByEmailAndVerified(email, true)
}

func (u UserRepository) Verify(user User) error {
	transaction, err := database.Connection.Begin()
	if err != nil {
		logrus.Errorf("failed creating transaction: %s", err)
		return err
	}

	sql, params, err := squirrel.Insert(verifiedTableName).SetMap(user.AsMap()).ToSql()
	if err != nil {
		logrus.Errorf("failed building INSERT query: %s", err)
		return err
	}

	rows, err := transaction.Exec(sql, params...)
	if err != nil {
		logrus.Errorf("failed executing INSERT query: %s", err)
		return err
	}
	if affected, err2 := rows.RowsAffected(); err2 != nil || affected < 1 {
		logrus.Errorf("not enough rows affected by INSERT: %s", err)
		if err != nil {
			return err
		} else {
			return errors.New("no rows inserted")
		}
	}

	sql, params, err = squirrel.Delete(unverifiedTableName).Where("id = ?", user.Id).ToSql()
	rows, err = transaction.Exec(sql, params...)
	if err != nil {
		logrus.Errorf("failed executing DELETE query: %s", err)
		txErr := transaction.Rollback()
		if txErr != nil {
			logrus.Errorf("Could not roll back transaction: %s", err)
			return txErr
		}
		return err
	}
	if affected, err2 := rows.RowsAffected(); err2 != nil || affected < 1 {
		logrus.Errorf("not enough rows affected by DELETE: %s", err)
		txErr := transaction.Rollback()
		if txErr != nil {
			logrus.Errorf("Could not roll back transaction: %s", err)
			return txErr
		}

		if err != nil {
			return err
		} else {
			return errors.New("no rows deleted")
		}
	}

	return transaction.Commit()
}

func (u UserRepository) Save(user User) error {
	transaction, err := database.Connection.Begin()
	if err != nil {
		logrus.Errorf("Could not save user: Could not create database transaction for save: %s", err)
		return err
	}

	var tableName string
	if user.Verified {
		tableName = verifiedTableName
	} else {
		tableName = unverifiedTableName
	}

	sql, params, err := squirrel.
		Insert(tableName).
		SetMap(user.AsMap()).
		ToSql()

	if user.Id > 0 {
		sql, params, err = squirrel.
			Update(tableName).
			SetMap(user.AsUpdateMap()).
			Where("id = ?", user.Id).
			ToSql()
	}

	if err != nil {
		logrus.Errorf("Could not save user: Could not build SQL query for save: %s", err)
		return err
	}

	logrus.Tracef("Executing SQL query: %s with params: %s", sql, params)

	result, err := transaction.Exec(sql, params...)
	if err != nil {
		logrus.Errorf("Could not execute SQL query for saving user: %s", err)
		_ = transaction.Rollback()
		return err
	}

	err = transaction.Commit()
	if err != nil {
		logrus.Errorf("Could not commit transaction for saving user: %s", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil || affected < 1 {
		logrus.Errorf("Could not save user: %s", err)
		return err
	}

	return nil
}

func (u UserRepository) FindByIdAndVerified(id int64, verified bool) (*User, error) {
	var usr *User
	cols := getCols(verified)

	tableName := verifiedTableName
	if !verified {
		tableName = unverifiedTableName
	}

	sql, params, err := squirrel.
		Select(cols...).
		From(tableName).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		logrus.Errorf("Could not build SQL query for fetching user: %s", err)
		return usr, err
	}

	rows, err := database.Connection.Queryx(sql, params...)
	if err != nil {
		logrus.Errorf("Could not run SQL query for fetching user: %s", err)
		logrus.Debugf("query: %s", sql)
		return usr, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Errorf("Could not close rows: %s", err)
		}
	}()

	if rows.Next() {
		usr = new(User)
		err = rows.StructScan(usr)
		usr.Verified = usr.VerificationToken == ""
		if err != nil {
			logrus.Errorf("Could not scan user into struct: %s", err)
			return usr, err
		} else {
			return usr, nil
		}
	} else {
		// no user found is not an error condition, it's expected - just return nil for user instead
		return nil, nil
	}
}

func getCols(verified bool) []string {
	cols := []string{
		"id",
		"email",
		"password_hash",
		"roles",
		"balance",
	}
	if !verified {
		cols = append(cols, "verification_token")
	}
	return cols
}

func (u UserRepository) FindVerified(id int64) (*User, error) {
	return u.FindByIdAndVerified(id, true)
}

func (u UserRepository) FindUnverified(id int64) (*User, error) {
	return u.FindByIdAndVerified(id, false)
}

func (u UserRepository) Find(id int64) (*User, error) {
	verifiedUser, err := u.FindVerified(id)
	if err != nil {
		return nil, err
	}
	if verifiedUser != nil {
		return verifiedUser, nil
	}

	return u.FindUnverified(id)
}

func (u UserRepository) FindAll() ([]User, error) {
	sql, params, err := squirrel.Select("id", "email", "balance", "password_hash", "roles").From(verifiedTableName).ToSql()

	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	err = database.Connection.Select(&users, sql, params...)

	return users, err
}

func (u UserRepository) FindByVerificationToken(token string) (*User, error) {
	sql, params, err := squirrel.Select(
		getCols(false)...,
	).
		From(unverifiedTableName).
		Where("verification_token = ?", token).
		ToSql()

	var usr *User

	rows, err := database.Connection.Queryx(sql, params...)
	if err != nil {
		logrus.Errorf("Could not run SQL query for fetching user: %s", err)
		return usr, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Errorf("Could not close rows: %s", err)
		}
	}()

	if rows.Next() {
		usr = new(User)
		err = rows.StructScan(usr)
		usr.Verified = usr.VerificationToken == ""
		if err != nil {
			logrus.Errorf("Could not scan user into struct: %s", err)
			return usr, err
		} else {
			return usr, nil
		}
	} else {
		// no user found is not an error condition, it's expected - just return nil for user instead
		return nil, nil
	}

}
