-- name: UpsertCoupons :exec
INSERT INTO "stripe"."coupons" (
    id, object, amount_off, created, currency, duration, duration_in_months, livemode, max_redemptions,
    metadata, name, percent_off, redeem_by, times_redeemed, valid, applies_to, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW()
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    amount_off = EXCLUDED.amount_off,
    created = EXCLUDED.created,
    currency = EXCLUDED.currency,
    duration = EXCLUDED.duration,
    duration_in_months = EXCLUDED.duration_in_months,
    livemode = EXCLUDED.livemode,
    max_redemptions = EXCLUDED.max_redemptions,
    metadata = EXCLUDED.metadata,
    name = EXCLUDED.name,
    percent_off = EXCLUDED.percent_off,
    redeem_by = EXCLUDED.redeem_by,
    times_redeemed = EXCLUDED.times_redeemed,
    valid = EXCLUDED.valid,
    applies_to = EXCLUDED.applies_to,
    updated_at = NOW();

-- name: DeleteCoupon :exec
DELETE FROM "stripe"."coupons" WHERE id = $1;

-- name: DeleteAllCoupons :exec
DELETE FROM "stripe"."coupons";

-- name: CouponExists :one
SELECT EXISTS (SELECT 1 FROM "stripe"."coupons" WHERE id = $1);