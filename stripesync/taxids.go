package stripesync

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
)

const excludedFieldCustomerTaxIDs = "customer.tax_ids"

func (o *StripeSync) taxIDsExcluded() bool {
	for _, field := range o.cfg.ExcludedFields {
		if field == excludedFieldCustomerTaxIDs {
			return true
		}
	}
	return false
}

// handleTaxIDUpdated relies on the event payload being complete, so unlike coupons it never has
// to fetch the tax ID back from Stripe.
func (o *StripeSync) handleTaxIDUpdated(ctx context.Context, taxID *stripe.TaxID) error {
	if o.taxIDsExcluded() {
		return nil
	}

	// Erroring here would abort the whole sync, and keep aborting it on every retry.
	if taxID.Customer == nil {
		log.Info().Str("tax_id", taxID.ID).Msg("Skipping tax id without a customer")
		return nil
	}

	err := o.ensureCustomerLoaded(ctx, taxID.Customer.ID)
	if err != nil {
		return err
	}

	return o.upsertTaxID(ctx, taxID, taxID.Customer.ID)
}

func (o *StripeSync) handleTaxIDDeleted(ctx context.Context, taxID *stripe.TaxID) error {
	if o.taxIDsExcluded() {
		return nil
	}

	return o.db.Q.DeleteCustomerTaxID(ctx, taxID.ID)
}

func (o *StripeSync) upsertTaxID(ctx context.Context, taxID *stripe.TaxID, customerID string) error {
	return o.db.Q.UpsertCustomerTaxID(ctx, postgres.UpsertCustomerTaxIDParams{
		ID:           taxID.ID,
		Object:       taxID.Object,
		Country:      utils.StringToNullString(taxID.Country),
		Created:      taxID.Created,
		Livemode:     taxID.Livemode,
		Type:         string(taxID.Type),
		Value:        taxID.Value,
		Verification: utils.MarshalToNullRawMessage(taxID.Verification),
		Customer:     customerID,
	})
}

// syncCustomerTaxIDs deletes any stored tax ID that is not in taxIDs, so only call it with a list
// Stripe returned in full, that is, from an expanded customer.
func (o *StripeSync) syncCustomerTaxIDs(ctx context.Context, customerID string, taxIDs []*stripe.TaxID) error {
	if o.taxIDsExcluded() {
		return nil
	}

	keepIDs := make([]string, 0, len(taxIDs))
	for _, taxID := range taxIDs {
		err := o.upsertTaxID(ctx, taxID, customerID)
		if err != nil {
			return fmt.Errorf("failed to upsert tax id %s: %w", taxID.ID, err)
		}
		keepIDs = append(keepIDs, taxID.ID)
	}

	return o.db.Q.DeleteCustomerTaxIDsOfCustomerExcept(ctx, postgres.DeleteCustomerTaxIDsOfCustomerExceptParams{
		Customer: customerID,
		KeepIds:  keepIDs,
	})
}
