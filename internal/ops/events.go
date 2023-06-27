package ops

import (
	"context"
	"encoding/json"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
)

func unmarshalEventData[T interface{}](e *stripe.Event) (*T, error) {
	var data T
	err := json.Unmarshal(e.Data.Raw, &data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal event data")
	}
	return &data, err
}

func HandleEvent(c context.Context, db *postgres.PostgresStore, e *stripe.Event) error {
	switch e.Type {
	case "customer.created", "customer.updated", "customer.deleted":
		customer, err := unmarshalEventData[stripe.Customer](e)
		if err == nil {
			return HandleCustomerUpdated(c, db, customer)
		}
		break
	case "product.created", "product.updated":
		product, err := unmarshalEventData[stripe.Product](e)
		if err == nil {
			return HandleProductUpdated(c, db, product)
		}
		break
	case "product.deleted":
		product, err := unmarshalEventData[stripe.Product](e)
		if err == nil {
			return HandleProductDeleted(c, db, product)
		}
		break
	case "subscription.created", "subscription.updated", "subscription.deleted":
		subscription, err := unmarshalEventData[stripe.Subscription](e)
		if err == nil {
			return HandleSubscriptionUpdated(c, db, subscription)
		}
		break
	case "price.created", "price.updated":
		price, err := unmarshalEventData[stripe.Price](e)
		if err == nil {
			return HandlePriceUpdated(c, db, price)
		}
		break
	case "price.deleted":
		price, err := unmarshalEventData[stripe.Price](e)
		if err == nil {
			return HandlePriceDeleted(c, db, price)
		}
		break
	case "coupon.created", "coupon.updated":
		coupon, err := unmarshalEventData[stripe.Coupon](e)
		if err == nil {
			return HandleCouponUpdated(c, db, coupon)
		}
		break
	case "coupon.deleted":
		coupon, err := unmarshalEventData[stripe.Coupon](e)
		if err == nil {
			return HandleCouponDeleted(c, db, coupon)
		}
		break
	case "customer.discount.updated", "customer.discount.deleted":
		discount, err := unmarshalEventData[stripe.Discount](e)
		if err == nil {
			return HandleSubscriptionDiscountUpdated(c, db, discount)
		}
		break
	}
	return nil
}
