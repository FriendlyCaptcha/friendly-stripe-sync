package stripesync

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

func (o *StripeSync) handleEvent(c context.Context, e *stripe.Event) error {
	switch e.Type {
	case "customer.created", "customer.updated", "customer.deleted":
		customer, err := unmarshalEventData[stripe.Customer](e)
		if err == nil {
			return o.handleCustomerUpdated(c, customer)
		}
	case "product.created", "product.updated":
		product, err := unmarshalEventData[stripe.Product](e)
		if err == nil {
			return o.handleProductUpdated(c, product)
		}
	case "product.deleted":
		product, err := unmarshalEventData[stripe.Product](e)
		if err == nil {
			return o.handleProductDeleted(c, product)
		}
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted", "customer.subscription.paused":
		subscription, err := unmarshalEventData[stripe.Subscription](e)
		if err == nil {
			return o.handleSubscriptionUpdated(c, subscription)
		}
	case "price.created", "price.updated":
		price, err := unmarshalEventData[stripe.Price](e)
		if err == nil {
			return o.handlePriceUpdated(c, price)
		}
	case "price.deleted":
		price, err := unmarshalEventData[stripe.Price](e)
		if err == nil {
			return o.handlePriceDeleted(c, price)
		}
	case "coupon.created", "coupon.updated":
		coupon, err := unmarshalEventData[stripe.Coupon](e)
		if err == nil {
			return o.handleCouponUpdated(c, coupon)
		}
	case "coupon.deleted":
		coupon, err := unmarshalEventData[stripe.Coupon](e)
		if err == nil {
			return o.handleCouponDeleted(c, coupon)
		}
	case "customer.discount.updated", "customer.discount.deleted":
		discount, err := unmarshalEventData[stripe.Discount](e)
		if err == nil {
			return o.handleSubscriptionDiscountUpdated(c, discount)
		}
	}
	return nil
}
