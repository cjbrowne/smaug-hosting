package transactions

import (
	"bitbucket.org/smaug-hosting/services/database"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"time"
)

const tableName = "pending_transactions"
const completedTableName = "completed_transactions"

type PendingTransactionRepository struct {
}

func (repository PendingTransactionRepository) Save(transaction PendingTransaction) {
	criticalLogger := logrus.WithField("severity", "CRITICAL")

	sql, params, err := squirrel.Insert(tableName).SetMap(map[string]interface{}{
		"user_id":     transaction.UserId,
		"amount":      transaction.Amount,
		"checkout_id": transaction.CheckoutId,
	}).ToSql()

	if err != nil {
		criticalLogger.Errorf("Could not generate SQL for saving pending transaction: %s", err)
		return
	}

	result, err := database.Connection.Exec(sql, params...)
	if err != nil {
		criticalLogger.Errorf("Could not save pending transaction: %s", err)
		return
	}

	if affected, err := result.RowsAffected(); affected != 1 || err != nil {
		criticalLogger.Errorf("Did not successfully save pending transaction: %s", err)
		return
	}
}

func (repository PendingTransactionRepository) FindByCheckoutId(checkoutId string) (*PendingTransaction, error) {
	sql, params, err := squirrel.Select("user_id", "amount", "checkout_id").From(tableName).Where("checkout_id = ?", checkoutId).ToSql()
	if err != nil {
		logrus.Errorf("Could not build query for fetching pending transaction: %s", err)
		return nil, err
	}

	row := database.Connection.QueryRowx(sql, params...)
	if row.Err() != nil {
		logrus.Errorf("Could not fetch pending transaction: %s", row.Err())
		return nil, row.Err()
	}

	pending := new(PendingTransaction)

	err = row.StructScan(pending)

	return pending, err
}

func (repository PendingTransactionRepository) MarkAsCompleted(transaction *PendingTransaction) error {
	tx, err := database.Connection.Beginx()
	if err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s", transaction.UserId, transaction.Amount, err)
		return err
	}

	sql, params, err := squirrel.Insert(completedTableName).SetMap(map[string]interface{}{
		"user_id":     transaction.UserId,
		"amount":      transaction.Amount,
		"checkout_id": transaction.CheckoutId,
		"completed":   time.Now(),
	}).ToSql()

	res, err := tx.Exec(sql, params...)

	if err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s", transaction.UserId, transaction.Amount, err)
		return err
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s (rows affected: %d)", transaction.UserId, transaction.Amount, err, rows)
		return err
	}

	sql, params, err = squirrel.Delete(tableName).Where("checkout_id = ?", transaction.CheckoutId).ToSql()
	if err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s", transaction.UserId, transaction.Amount, err)
		return err
	}

	res, err = tx.Exec(sql, params...)

	if err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s", transaction.UserId, transaction.Amount, err)
		_ = tx.Rollback()
		return err
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s", transaction.UserId, transaction.Amount, err)
		_ = tx.Rollback()
		if err == nil {
			return errors.New("rows_affected != 1")
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("Could not mark transaction as completed (user %d may get duplicate top-up balance of %d!!): %s", transaction.UserId, transaction.Amount, err)
		return err
	}

	return nil
}

func (repository PendingTransactionRepository) FindCompletedByCheckoutId(checkoutId string) (*PendingTransaction, error) {
	sql, params, err := squirrel.Select("user_id", "amount").From(completedTableName).Where("checkout_id = ?", checkoutId).ToSql()
	if err != nil {
		logrus.Errorf("Could not build query for fetching pending transaction: %s", err)
		return nil, err
	}

	row := database.Connection.QueryRowx(sql, params...)
	if row.Err() != nil {
		logrus.Errorf("Could not fetch pending transaction: %s", row.Err())
		return nil, row.Err()
	}

	pending := new(PendingTransaction)

	err = row.StructScan(pending)

	return pending, err
}
