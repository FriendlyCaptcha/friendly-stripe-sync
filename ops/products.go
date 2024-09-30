package ops

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/stripe/stripe-go/v74"
)

func (o *Ops) HandleProductUpdated(c context.Context, product *stripe.Product) error {
	return o.db.Q.UpsertProduct(c, postgres.UpsertProductParams{
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

func (o *Ops) HandleProductDeleted(c context.Context, product *stripe.Product) error {
	return o.db.Q.DeleteProduct(c, product.ID)
}

func (o *Ops) EnsureProductLoaded(c context.Context, productId string) error {
	exists, err := o.db.Q.ProductExists(c, productId)
	if err != nil {
		return err
	}
	if !exists {
		product, err := o.stripe.Products.Get(productId, nil)
		if err != nil {
			return err
		}
		err = o.HandleProductUpdated(c, product)
		if err != nil {
			return fmt.Errorf("failed to upsert product: %w", err)
		}
	}

	return nil
}
