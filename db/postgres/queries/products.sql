-- name: UpsertProduct :exec
INSERT INTO "stripe"."products" (
    id, object, active, description, metadata, name, created, images, livemode, package_dimensions,
    shippable, statement_descriptor, unit_label, updated, url, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW()
) ON CONFLICT (id) DO UPDATE SET 
    object = EXCLUDED.object,
    active = EXCLUDED.active,
    description = EXCLUDED.description,
    metadata = EXCLUDED.metadata,
    name = EXCLUDED.name,
    created = EXCLUDED.created,
    images = EXCLUDED.images,
    livemode = EXCLUDED.livemode,
    package_dimensions = EXCLUDED.package_dimensions,
    shippable = EXCLUDED.shippable,
    statement_descriptor = EXCLUDED.statement_descriptor,
    unit_label = EXCLUDED.unit_label,
    updated = EXCLUDED.updated,
    url = EXCLUDED.url,
    updated_at = NOW();

-- name: DeleteProduct :exec
DELETE FROM "stripe"."products" WHERE id = $1;

-- name: ProductExists :one
SELECT EXISTS (SELECT 1 FROM "stripe"."products" WHERE id = $1);

-- name: DeleteAllProducts :exec
DELETE FROM "stripe"."products";