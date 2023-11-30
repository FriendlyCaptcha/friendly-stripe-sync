package ops

import (
	"context"
	"database/sql"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/stripe/stripe-go/v74"
)

func HandleSubscriptionUpdated(c context.Context, db *postgres.PostgresStore, subscription *stripe.Subscription) error {
	err := EnsureCustomerLoaded(c, db, subscription.Customer.ID)
	if err != nil {
		return err
	}

	discountID := sql.NullString{}
	discountStart := sql.NullInt64{}
	discountEnd := sql.NullInt64{}
	discountCoupon := sql.NullString{}
	discountDeleted := sql.NullBool{}
	discountPromotionCode := sql.NullString{}
	if subscription.Discount != nil {
		discountID = utils.StringToNullString(subscription.Discount.ID)
		discountStart = utils.Int64ToNullInt64(subscription.Discount.Start)
		discountEnd = utils.Int64ToNullInt64(subscription.Discount.End)
		discountCoupon = utils.StringToNullString(subscription.Discount.Coupon.ID)
		discountDeleted = sql.NullBool{Bool: subscription.Discount.Deleted, Valid: true}
		if subscription.Discount.PromotionCode != nil {
			discountPromotionCode = utils.StringToNullString(subscription.Discount.PromotionCode.ID)
		}

		err := EnsureCouponLoaded(c, db, subscription.Discount.Coupon.ID)
		if err != nil {
			return err
		}
	}

	err = db.Q.UpsertSubscription(c, postgres.UpsertSubscriptionParams{
		ID:                            subscription.ID,
		Object:                        subscription.Object,
		CancelAtPeriodEnd:             subscription.CancelAtPeriodEnd,
		CurrentPeriodEnd:              subscription.CurrentPeriodEnd,
		CurrentPeriodStart:            subscription.CurrentPeriodStart,
		Metadata:                      utils.MarshalToNullRawMessage(subscription.Metadata),
		PendingUpdate:                 utils.MarshalToNullRawMessage(subscription.PendingUpdate),
		Status:                        string(subscription.Status),
		ApplicationFeePercent:         sql.NullFloat64{Float64: subscription.ApplicationFeePercent, Valid: subscription.ApplicationFeePercent != 0},
		BillingCycleAnchor:            subscription.BillingCycleAnchor,
		BillingThresholds:             utils.MarshalToNullRawMessage(subscription.BillingThresholds),
		CancelAt:                      utils.Int64ToNullInt64(subscription.CancelAt),
		CanceledAt:                    utils.Int64ToNullInt64(subscription.CanceledAt),
		CollectionMethod:              string(subscription.CollectionMethod),
		Created:                       subscription.Created,
		DaysUntilDue:                  utils.Int64ToNullInt64(subscription.DaysUntilDue),
		DefaultTaxRates:               utils.MarshalToNullRawMessage(subscription.DefaultTaxRates),
		EndedAt:                       utils.Int64ToNullInt64(subscription.EndedAt),
		Livemode:                      subscription.Livemode,
		NextPendingInvoiceItemInvoice: utils.Int64ToNullInt64(subscription.NextPendingInvoiceItemInvoice),
		PauseCollection:               utils.MarshalToNullRawMessage(subscription.PauseCollection),
		PendingInvoiceItemInterval:    utils.MarshalToNullRawMessage(subscription.PendingInvoiceItemInterval),
		StartDate:                     subscription.StartDate,
		TransferData:                  utils.MarshalToNullRawMessage(subscription.TransferData),
		TrialEnd:                      utils.Int64ToNullInt64(subscription.TrialEnd),
		TrialStart:                    utils.Int64ToNullInt64(subscription.TrialStart),
		DiscountID:                    discountID,
		DiscountStart:                 discountStart,
		DiscountEnd:                   discountEnd,
		DiscountDeleted:               discountDeleted,
		DiscountPromotionCode:         discountPromotionCode,
		DiscountCoupon:                discountCoupon,
		Customer:                      subscription.Customer.ID,
	})
	if err != nil {
		return err
	}

	for _, item := range subscription.Items.Data {
		err := EnsurePriceLoaded(c, db, item.Price.ID)
		if err != nil {
			return err
		}

		err = db.Q.UpsertSubscriptionItem(c, postgres.UpsertSubscriptionItemParams{
			ID:                item.ID,
			Object:            item.Object,
			BillingThresholds: utils.MarshalToNullRawMessage(item.BillingThresholds),
			Created:           item.Created,
			Metadata:          utils.MarshalToNullRawMessage(item.Metadata),
			Price:             item.Price.ID,
			Quantity:          item.Quantity,
			Subscription:      subscription.ID,
			TaxRates:          utils.MarshalToNullRawMessage(item.TaxRates),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleSubscriptionDiscountUpdated(c context.Context, db *postgres.PostgresStore, discount *stripe.Discount) error {
	err := EnsureCouponLoaded(c, db, discount.Coupon.ID)
	if err != nil {
		return err
	}

	var promotionCodeID sql.NullString
	if discount.PromotionCode != nil {
		promotionCodeID = sql.NullString{
			Valid:  true,
			String: discount.PromotionCode.ID,
		}
	}

	var couponID sql.NullString
	if discount.Coupon != nil {
		couponID = sql.NullString{
			Valid:  true,
			String: discount.Coupon.ID,
		}
	}

	return db.Q.UpdateSubscriptionDiscount(c, postgres.UpdateSubscriptionDiscountParams{
		DiscountID:            utils.StringToNullString(discount.ID),
		DiscountStart:         utils.Int64ToNullInt64(discount.Start),
		DiscountEnd:           utils.Int64ToNullInt64(discount.End),
		DiscountDeleted:       sql.NullBool{Bool: discount.Deleted, Valid: true},
		DiscountPromotionCode: promotionCodeID,
		DiscountCoupon:        couponID,
	})
}
