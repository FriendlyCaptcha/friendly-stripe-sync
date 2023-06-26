/* https://stripe.com/docs/api/subscriptions/object */

create table if not exists "stripe"."subscriptions" (
  "id" text primary key,
  "object" text not null,
  "cancel_at_period_end" boolean not null,
  "current_period_end" bigint not null,
  "current_period_start" bigint not null,
  "default_payment_method" text,
  "metadata" jsonb,
  "pending_setup_intent" text,
  "pending_update" jsonb,
  "status" text not null,
  "application_fee_percent" double precision,
  "billing_cycle_anchor" bigint not null,
  "billing_thresholds" jsonb,
  "cancel_at" bigint,
  "canceled_at" bigint,
  "collection_method" text not null,
  "created" bigint not null,
  "days_until_due" bigint,
  "default_source" text,
  "default_tax_rates" jsonb,
  "discount" jsonb,
  "ended_at" bigint,
  "livemode" boolean not null,
  "next_pending_invoice_item_invoice" bigint,
  "pause_collection" jsonb,
  "pending_invoice_item_interval" jsonb,
  "start_date" bigint not null,
  "transfer_data" jsonb,
  "trial_end" bigint,
  "trial_start" bigint,

  "schedule" text,
  "customer" text not null references "stripe"."customers" on delete cascade,

  "updated_at" timestamptz not null
);