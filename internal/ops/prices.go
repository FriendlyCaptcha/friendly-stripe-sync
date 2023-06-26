package ops

import (
	"context"
	"database/sql"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
)

func HandlePriceUpdated(c context.Context, db *postgres.PostgresStore, price *stripe.Price) error {
	err := EnsureProductLoaded(c, db, price.Product.ID)
	if err != nil {
		return err
	}
	return db.Q.UpsertPrice(c, postgres.UpsertPriceParams{
		ID:                price.ID,
		Object:            price.Object,
		Active:            price.Active,
		BillingScheme:     string(price.BillingScheme),
		Created:           price.Created,
		Currency:          string(price.Currency),
		Livemode:          price.Livemode,
		LookupKey:         utils.StringToNullString(price.LookupKey),
		Metadata:          utils.MarshalToNullRawMessage(price.Metadata),
		Nickname:          utils.StringToNullString(price.Nickname),
		Recurring:         utils.MarshalToNullRawMessage(price.Recurring),
		Type:              string(price.Type),
		UnitAmount:        sql.NullInt64{Int64: price.UnitAmount, Valid: price.UnitAmount != 0},
		TiersMode:         utils.StringToNullString(string(price.TiersMode)),
		TransformQuantity: utils.MarshalToNullRawMessage(price.TransformQuantity),
		UnitAmountDecimal: sql.NullFloat64{Float64: price.UnitAmountDecimal, Valid: price.UnitAmountDecimal != 0},
		Product:           price.Product.ID,
	})
}

func HandlePriceDeleted(c context.Context, db *postgres.PostgresStore, price *stripe.Price) error {
	return db.Q.DeleteProduct(c, price.ID)
}

func EnsurePriceLoaded(c context.Context, db *postgres.PostgresStore, priceID string) error {
	exists, err := db.Q.PriceExists(c, priceID)
	if err != nil {
		return err
	}
	if !exists {
		price, err := price.Get(priceID, nil)
		if err != nil {
			return err
		}
		HandlePriceUpdated(c, db, price)
	}

	return nil
}
