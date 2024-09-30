package stripesync

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/sqlc-dev/pqtype"
	"github.com/stripe/stripe-go/v74"
)

func (o *StripeSync) handleCustomerUpdated(c context.Context, customer *stripe.Customer) error {
	address := utils.MarshalToNullRawMessage(customer.Address)
	phone := utils.StringToNullString(customer.Phone)

	for _, field := range o.cfg.ExcludedFields {
		switch field {
		case "customer.address":
			address = pqtype.NullRawMessage{}
		case "customer.phone":
			phone = sql.NullString{}
		}
	}

	return o.db.Q.UpsertCustomer(c, postgres.UpsertCustomerParams{
		ID:                  customer.ID,
		Object:              customer.Object,
		Address:             address,
		Description:         utils.StringToNullString(customer.Description),
		Email:               utils.StringToNullString(customer.Email),
		Metadata:            utils.MarshalToNullRawMessage(customer.Metadata),
		Name:                utils.StringToNullString(customer.Name),
		Phone:               phone,
		Shipping:            utils.MarshalToNullRawMessage(customer.Shipping),
		Balance:             sql.NullInt64{Int64: customer.Balance, Valid: true},
		Created:             customer.Created,
		Currency:            utils.StringToNullString(string(customer.Currency)),
		Delinquent:          customer.Delinquent,
		Discount:            utils.MarshalToNullRawMessage(customer.Discount),
		InvoicePrefix:       customer.InvoicePrefix,
		InvoiceSettings:     utils.MarshalToNullRawMessage(customer.InvoiceSettings),
		Livemode:            customer.Livemode,
		NextInvoiceSequence: customer.NextInvoiceSequence,
		PreferredLocales:    customer.PreferredLocales,
		TaxExempt:           string(customer.TaxExempt),
		Deleted:             customer.Deleted,
	})
}

func (o *StripeSync) ensureCustomerLoaded(ctx context.Context, customerID string) error {
	exists, err := o.db.Q.CustomerExists(ctx, customerID)
	if err != nil {
		return err
	}
	if !exists {
		customer, err := o.stripe.Customers.Get(customerID, nil)
		if err != nil {
			return err
		}
		err = o.handleCustomerUpdated(ctx, customer)
		if err != nil {
			return fmt.Errorf("failed to upsert customer: %w", err)
		}
	}

	return nil
}
