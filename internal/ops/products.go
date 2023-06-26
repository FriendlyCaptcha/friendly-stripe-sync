package ops

import (
	"context"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/product"
)

func HandleProductUpdated(c context.Context, db *postgres.PostgresStore, product *stripe.Product) error {
	return db.Q.UpsertProduct(c, postgres.UpsertProductParams{
		ID:                  product.ID,
		Object:              product.Object,
		Active:              product.Active,
		Description:         utils.StringToNullString(product.Description),
		Metadata:            utils.MarshalToNullRawMessage(product.Metadata),
		Name:                product.Name,
		Created:             product.Created,
		Images:              product.Images,
		Livemode:            product.Livemode,
		PackageDimensions:   utils.MarshalToNullRawMessage(product.PackageDimensions),
		Shippable:           product.Shippable,
		StatementDescriptor: utils.StringToNullString(product.StatementDescriptor),
		UnitLabel:           utils.StringToNullString(product.UnitLabel),
		Updated:             product.Updated,
		Url:                 utils.StringToNullString(product.URL),
	})
}

func HandleProductDeleted(c context.Context, db *postgres.PostgresStore, product *stripe.Product) error {
	return db.Q.DeleteProduct(c, product.ID)
}

func EnsureProductLoaded(c context.Context, db *postgres.PostgresStore, productId string) error {
	exists, err := db.Q.ProductExists(c, productId)
	if err != nil {
		return err
	}
	if !exists {
		product, err := product.Get(productId, nil)
		if err != nil {
			return err
		}
		HandleProductUpdated(c, db, product)
	}

	return nil
}
