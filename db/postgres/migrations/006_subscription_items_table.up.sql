create table if not exists "stripe"."subscription_items" (
  "id" text primary key,
  "object" text not null,
  "billing_thresholds" jsonb,
  "created" bigint not null,
  "metadata" jsonb,
  "quantity" bigint not null,
  "price" text not null references "stripe"."prices" on delete cascade,
  "subscription" text not null references "stripe"."subscriptions" on delete cascade,
  "tax_rates" jsonb
);