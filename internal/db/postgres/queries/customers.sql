-- name: UpsertCustomer :exec
INSERT INTO "stripe"."customers" (
    id, object, address, description, email, metadata, name, phone, shipping, balance, created, 
    currency, delinquent, discount, invoice_prefix, invoice_settings, livemode, 
    next_invoice_sequence, preferred_locales, tax_exempt, deleted, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, 
    $20, $21, NOW()
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    address = EXCLUDED.address,
    description = EXCLUDED.description,
    email = EXCLUDED.email,
    metadata = EXCLUDED.metadata,
    name = EXCLUDED.name,
    phone = EXCLUDED.phone,
    shipping = EXCLUDED.shipping,
    balance = EXCLUDED.balance,
    created = EXCLUDED.created,
    currency = EXCLUDED.currency,
    delinquent = EXCLUDED.delinquent,
    discount = EXCLUDED.discount,
    invoice_prefix = EXCLUDED.invoice_prefix,
    invoice_settings = EXCLUDED.invoice_settings,
    livemode = EXCLUDED.livemode,
    next_invoice_sequence = EXCLUDED.next_invoice_sequence,
    preferred_locales = EXCLUDED.preferred_locales,
    tax_exempt = EXCLUDED.tax_exempt,
    deleted = EXCLUDED.deleted,
    updated_at = NOW();

-- name: DeleteAllCustomers :exec
DELETE FROM "stripe"."customers";

-- name: CustomerExists :one
SELECT EXISTS (SELECT 1 FROM "stripe"."customers" WHERE id = $1);