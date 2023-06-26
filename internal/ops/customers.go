package ops

import (
	"context"
	"database/sql"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/tabbed/pqtype"
)

func HandleCustomerUpdated(c context.Context, db *postgres.PostgresStore, customer *stripe.Customer) error {
	address := utils.MarshalToNullRawMessage(customer.Address)
	phone := utils.StringToNullString(customer.Phone)

	for _, field := range viper.GetStringSlice("stripe_sync.excluded_fields") {
		switch field {
		case "customer.address":
			address = pqtype.NullRawMessage{}
		case "customer.phone":
			phone = sql.NullString{}
		}
	}

	return db.Q.UpsertCustomer(c, postgres.UpsertCustomerParams{
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

func EnsureCustomerLoaded(c context.Context, db *postgres.PostgresStore, customerID string) error {
	exists, err := db.Q.CustomerExists(c, customerID)
	if err != nil {
		return err
	}
	if !exists {
		customer, err := customer.Get(customerID, nil)
		if err != nil {
			return err
		}
		HandleCustomerUpdated(c, db, customer)
	}

	return nil
}
