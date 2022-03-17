package transactions

type PendingTransaction struct {
	UserId     int64 `db:"user_id"`
	Amount     int64 `db:"amount"`
	CheckoutId string `db:"checkout_id"`
}
