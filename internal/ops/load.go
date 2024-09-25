package ops

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"golang.org/x/sync/errgroup"
)

func (o *Ops) InitialLoad(ctx context.Context, purge bool) error {
	if purge {
		log.Info().Msgf("Deleting all existing data from database")

		err := o.db.Q.DeleteCurrentSyncState(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete current sync state")
			return err
		}

		err = o.db.Q.DeleteAllCustomers(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all customers")
			return err
		}

		err = o.db.Q.DeleteAllProducts(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all products")
			return err
		}

		err = o.db.Q.DeleteAllPrices(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete all prices")
			return err
		}

		err = o.db.Q.DeleteAllSubscriptions(ctx)
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
		return o.loadCustomers(ctx)
	})
	g.Go(func() error {
		return o.loadProducts(ctx)
	})
	g.Go(func() error {
		return o.loadPrices(ctx)
	})
	g.Go(func() error {
		return o.loadSubscriptions(ctx)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	log.Info().Msgf("Finished loading data from Stripe")
	err := o.db.Q.SetSyncState(context.Background(), startedAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set sync state")
		return err
	}

	return nil
}

func (o *Ops) loadCustomers(c context.Context) error {
	customers := o.stripe.Customers.List(&stripe.CustomerListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for customers.Next() {
		cus := customers.Customer()
		err := o.HandleCustomerUpdated(c, cus)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded customer")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "customer").Msg("Finished loading customers")
	return nil
}

func (o *Ops) loadProducts(c context.Context) error {
	products := o.stripe.Products.List(&stripe.ProductListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for products.Next() {
		p := products.Product()
		err := o.HandleProductUpdated(c, p)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded product")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "product").Msg("Finished loading products")
	return nil
}

func (o *Ops) loadPrices(c context.Context) error {
	prices := o.stripe.Prices.List(&stripe.PriceListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for prices.Next() {
		p := prices.Price()
		err := o.HandlePriceUpdated(c, p)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded price")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "price").Msg("Finished loading prices")
	return nil
}

func (o *Ops) loadSubscriptions(c context.Context) error {
	subscriptions := o.stripe.Subscriptions.List(&stripe.SubscriptionListParams{ListParams: stripe.ListParams{Limit: stripe.Int64(100)}})
	count := 0
	for subscriptions.Next() {
		s := subscriptions.Subscription()
		err := o.HandleSubscriptionUpdated(c, s)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle loaded subscription")
			return err
		}
		count += 1
	}
	log.Debug().Int("count", count).Str("entity_type", "subscription").Msg("Finished loading subscriptions")
	return nil
}
