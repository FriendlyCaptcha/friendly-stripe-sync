-- name: UpsertCustomerTaxID :exec
INSERT INTO "stripe"."customer_tax_ids" (
    id, object, country, created, livemode, type, value, verification, customer, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    country = EXCLUDED.country,
    created = EXCLUDED.created,
    livemode = EXCLUDED.livemode,
    type = EXCLUDED.type,
    value = EXCLUDED.value,
    verification = EXCLUDED.verification,
    customer = EXCLUDED.customer,
    updated_at = NOW();

-- name: DeleteCustomerTaxID :exec
DELETE FROM "stripe"."customer_tax_ids" WHERE id = $1;

-- name: DeleteCustomerTaxIDsOfCustomer :exec
DELETE FROM "stripe"."customer_tax_ids" WHERE customer = $1;

-- name: DeleteCustomerTaxIDsOfCustomerExcept :exec
DELETE FROM "stripe"."customer_tax_ids" WHERE customer = $1 AND id <> ALL(@keep_ids::text[]);

-- name: DeleteAllCustomerTaxIDs :exec
DELETE FROM "stripe"."customer_tax_ids";
