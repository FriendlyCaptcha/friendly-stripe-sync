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
		CancelAt:                      sql.NullInt64{Int64: subscription.CancelAt, Valid: subscription.CancelAt != 0},
		CanceledAt:                    sql.NullInt64{Int64: subscription.CanceledAt, Valid: subscription.CanceledAt != 0},
		CollectionMethod:              string(subscription.CollectionMethod),
		Created:                       subscription.Created,
		DaysUntilDue:                  sql.NullInt64{Int64: subscription.DaysUntilDue, Valid: subscription.DaysUntilDue != 0},
		DefaultTaxRates:               utils.MarshalToNullRawMessage(subscription.DefaultTaxRates),
		Discount:                      utils.MarshalToNullRawMessage(subscription.Discount),
		EndedAt:                       sql.NullInt64{Int64: subscription.EndedAt, Valid: subscription.EndedAt != 0},
		Livemode:                      subscription.Livemode,
		NextPendingInvoiceItemInvoice: sql.NullInt64{Int64: subscription.NextPendingInvoiceItemInvoice, Valid: subscription.NextPendingInvoiceItemInvoice != 0},
		PauseCollection:               utils.MarshalToNullRawMessage(subscription.PauseCollection),
		PendingInvoiceItemInterval:    utils.MarshalToNullRawMessage(subscription.PendingInvoiceItemInterval),
		StartDate:                     subscription.StartDate,
		TransferData:                  utils.MarshalToNullRawMessage(subscription.TransferData),
		TrialEnd:                      sql.NullInt64{Int64: subscription.TrialEnd, Valid: subscription.TrialEnd != 0},
		TrialStart:                    sql.NullInt64{Int64: subscription.TrialStart, Valid: subscription.TrialStart != 0},
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
