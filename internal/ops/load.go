package ops

import (
	"context"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
	"github.com/stripe/stripe-go/v74/subscription"
	"golang.org/x/sync/errgroup"
)

func InititalLoad(db *postgres.PostgresStore) error {
	ctx := context.Background()

	if viper.GetBool("purge") {
		log.Info().Msgf("Deleting all existing data from database")

		err := db.Q.DeleteCurrentSyncState(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete current sync state")
			return err
		}

		err = db.Q.DeleteAllCustomers(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all customers")
			return err
		}

		err = db.Q.DeleteAllProducts(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all products")
			return err
		}

		err = db.Q.DeleteAllPrices(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all prices")
			return err
		}

		err = db.Q.DeleteAllSubscriptions(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all subscriptions")
			return err
		}

		log.Info().Msgf("Finished deleting all existing data from database")
	}

	startedAt := time.Now().UTC().Unix()

	log.Info().Msgf("Starting to load data from Stripe")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return loadCustomers(ctx, db)
	})
	g.Go(func() error {
		return loadProducts(ctx, db)
	})
	g.Go(func() error {
		return loadPrices(ctx, db)
	})
	g.Go(func() error {
		return loadSubscriptions(ctx, db)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	log.Info().Msgf("Finished loading data from Stripe")
	err := db.Q.SetSyncState(context.Background(), startedAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set sync state")
		return err
	}

	return nil
}

func loadCustomers(c context.Context, db *postgres.PostgresStore) error {
	customers := customer.List(&stripe.CustomerListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for customers.Next() {
		cus := customers.Customer()
		err := HandleCustomerUpdated(c, db, cus)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded customer")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "customer").Msg("Finished loading customers")
	return nil
}

func loadProducts(c context.Context, db *postgres.PostgresStore) error {
	products := product.List(&stripe.ProductListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for products.Next() {
		p := products.Product()
		err := HandleProductUpdated(c, db, p)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded product")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "product").Msg("Finished loading products")
	return nil
}

func loadPrices(c context.Context, db *postgres.PostgresStore) error {
	prices := price.List(&stripe.PriceListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for prices.Next() {
		p := prices.Price()
		err := HandlePriceUpdated(c, db, p)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded price")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "price").Msg("Finished loading prices")
	return nil
}

func loadSubscriptions(c context.Context, db *postgres.PostgresStore) error {
	subscriptions := subscription.List(&stripe.SubscriptionListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for subscriptions.Next() {
		s := subscriptions.Subscription()
		err := HandleSubscriptionUpdated(c, db, s)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded subscription")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "subscription").Msg("Finished loading subscriptions")
	return nil
}
