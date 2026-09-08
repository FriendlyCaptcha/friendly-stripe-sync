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

	// Stripe drops the tax IDs of a deleted customer, so storing one would contradict
	// the cleanup that handleCustomerUpdated does for the same customer.
	deleted, err := o.db.Q.CustomerIsDeleted(ctx, taxID.Customer.ID)
	if err != nil {
		return err
	}
	if deleted {
		log.Info().Str("tax_id", taxID.ID).Msg("Skipping tax id of a deleted customer")
		return nil
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
func (o *StripeSync) syncCustomerTaxIDs(ctx context.Context, customerID string, taxIDs *stripe.TaxIDList) error {
	if o.taxIDsExcluded() {
		return nil
	}

	keepIDs := make([]string, 0, len(taxIDs.Data))
	for _, taxID := range taxIDs.Data {
		err := o.upsertTaxID(ctx, taxID, customerID)
		if err != nil {
			return fmt.Errorf("failed to upsert tax id %s: %w", taxID.ID, err)
		}
		keepIDs = append(keepIDs, taxID.ID)
	}

	// Expansion returns only the first page. Deleting what Stripe left out would drop real tax IDs,
	// so leave the rest alone and let the customer.tax_id.* events correct them.
	if taxIDs.HasMore {
		log.Warn().Str("customer_id", customerID).Int("loaded", len(keepIDs)).
			Msg("Customer has more tax ids than Stripe expands, skipping removal of the rest")
		return nil
	}

	return o.db.Q.DeleteCustomerTaxIDsOfCustomerExcept(ctx, postgres.DeleteCustomerTaxIDsOfCustomerExceptParams{
		Customer: customerID,
		KeepIds:  keepIDs,
	})
}
