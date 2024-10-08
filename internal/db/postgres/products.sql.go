// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: products.sql

package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const deleteAllProducts = `-- name: DeleteAllProducts :exec
DELETE FROM "stripe"."products"
`

func (q *Queries) DeleteAllProducts(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllProducts)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM "stripe"."products" WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const productExists = `-- name: ProductExists :one
SELECT EXISTS (SELECT 1 FROM "stripe"."products" WHERE id = $1)
`

func (q *Queries) ProductExists(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRowContext(ctx, productExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const upsertProduct = `-- name: UpsertProduct :exec
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
    updated_at = NOW()
`

type UpsertProductParams struct {
	ID                  string
	Object              string
	Active              bool
	Description         sql.NullString
	Metadata            pqtype.NullRawMessage
	Name                string
	Created             int64
	Images              []string
	Livemode            bool
	PackageDimensions   pqtype.NullRawMessage
	Shippable           bool
	StatementDescriptor sql.NullString
	UnitLabel           sql.NullString
	Updated             int64
	Url                 sql.NullString
}

func (q *Queries) UpsertProduct(ctx context.Context, arg UpsertProductParams) error {
	_, err := q.db.ExecContext(ctx, upsertProduct,
		arg.ID,
		arg.Object,
		arg.Active,
		arg.Description,
		arg.Metadata,
		arg.Name,
		arg.Created,
		pq.Array(arg.Images),
		arg.Livemode,
		arg.PackageDimensions,
		arg.Shippable,
		arg.StatementDescriptor,
		arg.UnitLabel,
		arg.Updated,
		arg.Url,
	)
	return err
}
