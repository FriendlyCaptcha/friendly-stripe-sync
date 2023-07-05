# friendly-stripe-sync

Load and periodically synchronize data from Stripe to a Postgres database.

## How it works

1. Creates a schema called `stripe` in the postgres database with all the needed tables
2. Loads the initial data from Stripe and populates the tables
3. Periodically calls the `/events` endpoint of Stripe to load changes and update the database

### Supported Data Types

friendly-stripe-sync currently supports products, prices, customers, subscriptions, subscription items, and coupons.  
The data is updated using the following events:

- [x] customer.created
- [x] customer.updated
- [x] customer.deleted
- [x] product.created
- [x] product.updated
- [x] product.deleted
- [x] subscription.created
- [x] subscription.updated
- [x] subscription.deleted
- [x] price.created
- [x] price.updated
- [x] price.deleted

## Installation

```shell
# Go >= 1.17
go install github.com/FriendlyCaptcha/friendly-stripe-sync@latest
# Go < 1.17:
go install github.com/FriendlyCaptcha/friendly-stripe-sync
```

## Usage

### Setup Database

Before running the `friendly-stripe-sync` service you will want to migrate your database to create the necessary tables.

```shell
# Help on all migration-related commands for Postgres
friendly-stripe-sync migrate postgres --help

# List all available migrations
friendly-stripe-sync migrate postgres list

# List the current migration version
friendly-stripe-sync migrate postgres version

# Run all migrations to get up to the latest version
friendly-stripe-sync migrate postgres up
```

### Load the existing Stripe data

To load the initial dataset (required if you have data older than 30 days in your Stripe account).
2

```shell
# --purge will delete all existing data before loading the data from Stripe
friendly-stripe-sync load [--purge]
```

### Synchronize once

To load recent changes since the initial load you can synchronize once. This will work for changes that are up to 30 days old.

```shell
friendly-stripe-sync sync
```

### Synchronize and watch

If you want to keep your database up to date you should consider running the synchronization periodically.  
This will load the initial data if it hasn't been loaded already and will then synchronize periodically (see configuration).

```shell
# -i defines how often in seconds it will load changes from Stripe
friendly-stripe-sync watch [-i 60]
```

## Configuration

The app looks for a configuration file called `.friendly-stripe-sync.yml` in the working directory.

```yaml
# Sets the log devel to debug and adds more info to the log entires
debug: false
# Makes the log entries human readable instead of JSON
development: false

stripe:
  # You Stripe API key
  api_key: "..."

# Database configuration
postgres:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: ""
  dbname: "friendlystripe"
  sslmode: "disable"

# Exclude certain fields from being stored in the database. By default no fields are excluded.
stripe_sync:
  excluded_fields:
    # Currently only the following fields can be excluded:
    - "customer.address"
    - "customer.phone"
```

You can also set config fields using environment variables. The app will look for environment variables starting with `FSS_`. For example `FSS_POSTGRES__PASSWORD=1234567890` will overwrite the `postgres.password` field in the YAML config.
