package ops

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/stripe/stripe-go/v74"
)

func (o *StripeSync) HandlePriceUpdated(c context.Context, price *stripe.Price) error {
	err := o.EnsureProductLoaded(c, price.Product.ID)
	if err != nil {
		return err
	}
	return o.db.Q.UpsertPrice(c, postgres.UpsertPriceParams{
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

func (o *StripeSync) HandlePriceDeleted(c context.Context, price *stripe.Price) error {
	return o.db.Q.DeleteProduct(c, price.ID)
}

func (o *StripeSync) EnsurePriceLoaded(c context.Context, priceID string) error {
	exists, err := o.db.Q.PriceExists(c, priceID)
	if err != nil {
		return err
	}
	if !exists {
		price, err := o.stripe.Prices.Get(priceID, nil)
		if err != nil {
			return err
		}
		err = o.HandlePriceUpdated(c, price)
		if err != nil {
			return fmt.Errorf("failed to upsert price: %w", err)
		}
	}

	return nil
}
