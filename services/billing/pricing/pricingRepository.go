package pricing

import (
	"bitbucket.org/smaug-hosting/services/database"
	"github.com/Masterminds/squirrel"
)

//noinspection GoNameStartsWithPackageName
type PricingRepository struct{}

const tableName = "prices"

func (pr PricingRepository) FindPriceBySoftwareAndTier(software string, tier int) (Price, error) {
	var price Price

	sql, params, err := squirrel.
		Select("amount", "software", "tier").
		From(tableName).
		Where("software = ? AND tier = ?", software, tier).
		ToSql()

	if err != nil {
		return price, err
	}

	row := database.Connection.QueryRowx(sql, params...)

	err = row.StructScan(&price)

	return price, err
}
