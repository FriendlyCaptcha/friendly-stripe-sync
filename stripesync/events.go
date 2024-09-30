package ops

import (
	"context"
	"encoding/json"

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

func (o *StripeSync) HandleEvent(c context.Context, e *stripe.Event) error {
	switch e.Type {
	case "customer.created", "customer.updated", "customer.deleted":
		customer, err := unmarshalEventData[stripe.Customer](e)
		if err == nil {
			return o.HandleCustomerUpdated(c, customer)
		}
		break
	case "product.created", "product.updated":
		product, err := unmarshalEventData[stripe.Product](e)
		if err == nil {
			return o.HandleProductUpdated(c, product)
		}
		break
	case "product.deleted":
		product, err := unmarshalEventData[stripe.Product](e)
		if err == nil {
			return o.HandleProductDeleted(c, product)
		}
		break
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted", "customer.subscription.paused":
		subscription, err := unmarshalEventData[stripe.Subscription](e)
		if err == nil {
			return o.HandleSubscriptionUpdated(c, subscription)
		}
		break
	case "price.created", "price.updated":
		price, err := unmarshalEventData[stripe.Price](e)
		if err == nil {
			return o.HandlePriceUpdated(c, price)
		}
		break
	case "price.deleted":
		price, err := unmarshalEventData[stripe.Price](e)
		if err == nil {
			return o.HandlePriceDeleted(c, price)
		}
		break
	case "coupon.created", "coupon.updated":
		coupon, err := unmarshalEventData[stripe.Coupon](e)
		if err == nil {
			return o.HandleCouponUpdated(c, coupon)
		}
		break
	case "coupon.deleted":
		coupon, err := unmarshalEventData[stripe.Coupon](e)
		if err == nil {
			return o.HandleCouponDeleted(c, coupon)
		}
		break
	case "customer.discount.updated", "customer.discount.deleted":
		discount, err := unmarshalEventData[stripe.Discount](e)
		if err == nil {
			return o.HandleSubscriptionDiscountUpdated(c, discount)
		}
		break
	}
	return nil
}
