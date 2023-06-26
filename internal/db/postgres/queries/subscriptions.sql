-- name: UpsertSubscription :exec
INSERT INTO "stripe"."subscriptions" (
    id, object, cancel_at_period_end, current_period_end, current_period_start,
    metadata, pending_update, status, application_fee_percent,
    billing_cycle_anchor, billing_thresholds, cancel_at, canceled_at, collection_method, created,
    days_until_due, default_tax_rates, discount, ended_at, livemode,
    next_pending_invoice_item_invoice, pause_collection, pending_invoice_item_interval, start_date,
    transfer_data, trial_end, trial_start, customer, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21,
    $22, $23, $24, $25, $26, $27, $28, NOW()
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
    discount = EXCLUDED.discount,
    ended_at = EXCLUDED.ended_at,
    livemode = EXCLUDED.livemode,
    next_pending_invoice_item_invoice = EXCLUDED.next_pending_invoice_item_invoice,
    pause_collection = EXCLUDED.pause_collection,
    pending_invoice_item_interval = EXCLUDED.pending_invoice_item_interval,
    start_date = EXCLUDED.start_date,
    transfer_data = EXCLUDED.transfer_data,
    trial_end = EXCLUDED.trial_end,
    trial_start = EXCLUDED.trial_start,
    customer = EXCLUDED.customer,
    updated_at = NOW();

-- name: UpsertSubscriptionItem :exec
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
    tax_rates = EXCLUDED.tax_rates;

-- name: DeleteAllSubscriptions :exec
DELETE FROM "stripe"."subscriptions";