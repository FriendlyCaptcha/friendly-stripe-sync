-- name: UpsertPrice :exec
INSERT INTO "stripe"."prices" (
    id, object, active, billing_scheme, created, currency, livemode, lookup_key, metadata,
    nickname, recurring, type, unit_amount, tiers_mode, transform_quantity, unit_amount_decimal,
    product, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW()
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    active = EXCLUDED.active,
    billing_scheme = EXCLUDED.billing_scheme,
    created = EXCLUDED.created,
    currency = EXCLUDED.currency,
    livemode = EXCLUDED.livemode,
    lookup_key = EXCLUDED.lookup_key,
    metadata = EXCLUDED.metadata,
    nickname = EXCLUDED.nickname,
    recurring = EXCLUDED.recurring,
    type = EXCLUDED.type,
    unit_amount = EXCLUDED.unit_amount,
    tiers_mode = EXCLUDED.tiers_mode,
    transform_quantity = EXCLUDED.transform_quantity,
    unit_amount_decimal = EXCLUDED.unit_amount_decimal,
    product = EXCLUDED.product,
    updated_at = NOW();

-- name: DeletePrice :exec
DELETE FROM "stripe"."prices" WHERE id = $1;

-- name: DeleteAllPrices :exec
DELETE FROM "stripe"."prices";

-- name: PriceExists :one
SELECT EXISTS (SELECT 1 FROM "stripe"."prices" WHERE id = $1);