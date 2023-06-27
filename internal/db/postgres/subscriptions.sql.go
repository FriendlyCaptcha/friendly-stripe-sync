// Code generated by sqlc. DO NOT EDIT.
// source: subscriptions.sql

package postgres

import (
	"context"
	"database/sql"

	"github.com/tabbed/pqtype"
)

const deleteAllSubscriptions = `-- name: DeleteAllSubscriptions :exec
DELETE FROM "stripe"."subscriptions"
`

func (q *Queries) DeleteAllSubscriptions(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllSubscriptions)
	return err
}

const updateSubscriptionDiscount = `-- name: UpdateSubscriptionDiscount :exec
UPDATE "stripe"."subscriptions" SET
    discount_id = $1,
    discount_start = $2,
    discount_end = $3,
    discount_coupon = $4,
    discount_deleted = $5,
    discount_promotion_code = $6
WHERE discount_id = $1
`

type UpdateSubscriptionDiscountParams struct {
	DiscountID            sql.NullString
	DiscountStart         sql.NullInt64
	DiscountEnd           sql.NullInt64
	DiscountCoupon        sql.NullString
	DiscountDeleted       sql.NullBool
	DiscountPromotionCode sql.NullString
}

func (q *Queries) UpdateSubscriptionDiscount(ctx context.Context, arg UpdateSubscriptionDiscountParams) error {
	_, err := q.db.ExecContext(ctx, updateSubscriptionDiscount,
		arg.DiscountID,
		arg.DiscountStart,
		arg.DiscountEnd,
		arg.DiscountCoupon,
		arg.DiscountDeleted,
		arg.DiscountPromotionCode,
	)
	return err
}

const upsertSubscription = `-- name: UpsertSubscription :exec
INSERT INTO "stripe"."subscriptions" (
    id, object, cancel_at_period_end, current_period_end, current_period_start,
    metadata, pending_update, status, application_fee_percent,
    billing_cycle_anchor, billing_thresholds, cancel_at, canceled_at, collection_method, created,
    days_until_due, default_tax_rates, ended_at, livemode,
    next_pending_invoice_item_invoice, pause_collection, pending_invoice_item_interval, start_date,
    transfer_data, trial_end, trial_start, discount_id, discount_start, discount_end, discount_coupon, 
    discount_deleted, discount_promotion_code, customer, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21,
    $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, NOW()
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    cancel_at_period_end = EXCLUDED.cancel_at_period_end,
    current_period_end = EXCLUDED.current_period_end,
    current_period_start = EXCLUDED.current_period_start,
    metadata = EXCLUDED.metadata,
    pending_update = EXCLUDED.pending_update,
    status = EXCLUDED.status,
    application_fee_percent = EXCLUDED.application_fee_percent,
    billing_cycle_anchor = EXCLUDED.billing_cycle_anchor,
    billing_thresholds = EXCLUDED.billing_thresholds,
    cancel_at = EXCLUDED.cancel_at,
    canceled_at = EXCLUDED.canceled_at,
    collection_method = EXCLUDED.collection_method,
    created = EXCLUDED.created,
    days_until_due = EXCLUDED.days_until_due,
    default_tax_rates = EXCLUDED.default_tax_rates,
    ended_at = EXCLUDED.ended_at,
    livemode = EXCLUDED.livemode,
    next_pending_invoice_item_invoice = EXCLUDED.next_pending_invoice_item_invoice,
    pause_collection = EXCLUDED.pause_collection,
    pending_invoice_item_interval = EXCLUDED.pending_invoice_item_interval,
    start_date = EXCLUDED.start_date,
    transfer_data = EXCLUDED.transfer_data,
    trial_end = EXCLUDED.trial_end,
    trial_start = EXCLUDED.trial_start,
    discount_id = EXCLUDED.discount_id,
    discount_start = EXCLUDED.discount_start,
    discount_end = EXCLUDED.discount_end,
    discount_coupon = EXCLUDED.discount_coupon,
    discount_deleted = EXCLUDED.discount_deleted,
    discount_promotion_code = EXCLUDED.discount_promotion_code,
    customer = EXCLUDED.customer,
    updated_at = NOW()
`

type UpsertSubscriptionParams struct {
	ID                            string
	Object                        string
	CancelAtPeriodEnd             bool
	CurrentPeriodEnd              int64
	CurrentPeriodStart            int64
	Metadata                      pqtype.NullRawMessage
	PendingUpdate                 pqtype.NullRawMessage
	Status                        string
	ApplicationFeePercent         sql.NullFloat64
	BillingCycleAnchor            int64
	BillingThresholds             pqtype.NullRawMessage
	CancelAt                      sql.NullInt64
	CanceledAt                    sql.NullInt64
	CollectionMethod              string
	Created                       int64
	DaysUntilDue                  sql.NullInt64
	DefaultTaxRates               pqtype.NullRawMessage
	EndedAt                       sql.NullInt64
	Livemode                      bool
	NextPendingInvoiceItemInvoice sql.NullInt64
	PauseCollection               pqtype.NullRawMessage
	PendingInvoiceItemInterval    pqtype.NullRawMessage
	StartDate                     int64
	TransferData                  pqtype.NullRawMessage
	TrialEnd                      sql.NullInt64
	TrialStart                    sql.NullInt64
	DiscountID                    sql.NullString
	DiscountStart                 sql.NullInt64
	DiscountEnd                   sql.NullInt64
	DiscountCoupon                sql.NullString
	DiscountDeleted               sql.NullBool
	DiscountPromotionCode         sql.NullString
	Customer                      string
}

func (q *Queries) UpsertSubscription(ctx context.Context, arg UpsertSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, upsertSubscription,
		arg.ID,
		arg.Object,
		arg.CancelAtPeriodEnd,
		arg.CurrentPeriodEnd,
		arg.CurrentPeriodStart,
		arg.Metadata,
		arg.PendingUpdate,
		arg.Status,
		arg.ApplicationFeePercent,
		arg.BillingCycleAnchor,
		arg.BillingThresholds,
		arg.CancelAt,
		arg.CanceledAt,
		arg.CollectionMethod,
		arg.Created,
		arg.DaysUntilDue,
		arg.DefaultTaxRates,
		arg.EndedAt,
		arg.Livemode,
		arg.NextPendingInvoiceItemInvoice,
		arg.PauseCollection,
		arg.PendingInvoiceItemInterval,
		arg.StartDate,
		arg.TransferData,
		arg.TrialEnd,
		arg.TrialStart,
		arg.DiscountID,
		arg.DiscountStart,
		arg.DiscountEnd,
		arg.DiscountCoupon,
		arg.DiscountDeleted,
		arg.DiscountPromotionCode,
		arg.Customer,
	)
	return err
}

const upsertSubscriptionItem = `-- name: UpsertSubscriptionItem :exec
INSERT INTO "stripe"."subscription_items" (
    id, object, billing_thresholds, created, metadata, price, quantity, subscription, tax_rates
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    billing_thresholds = EXCLUDED.billing_thresholds,
    created = EXCLUDED.created,
    metadata = EXCLUDED.metadata,
    price = EXCLUDED.price,
    quantity = EXCLUDED.quantity,
    subscription = EXCLUDED.subscription,
    tax_rates = EXCLUDED.tax_rates
`

type UpsertSubscriptionItemParams struct {
	ID                string
	Object            string
	BillingThresholds pqtype.NullRawMessage
	Created           int64
	Metadata          pqtype.NullRawMessage
	Price             string
	Quantity          int64
	Subscription      string
	TaxRates          pqtype.NullRawMessage
}

func (q *Queries) UpsertSubscriptionItem(ctx context.Context, arg UpsertSubscriptionItemParams) error {
	_, err := q.db.ExecContext(ctx, upsertSubscriptionItem,
		arg.ID,
		arg.Object,
		arg.BillingThresholds,
		arg.Created,
		arg.Metadata,
		arg.Price,
		arg.Quantity,
		arg.Subscription,
		arg.TaxRates,
	)
	return err
}
