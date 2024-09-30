/* https://stripe.com/docs/api/prices/object */


create table if not exists "stripe"."prices" (
  "id" text primary key,
  "object" text not null,
  "active" boolean not null,
  "currency" text not null,
  "metadata" jsonb,
  "nickname" text,
  "recurring" jsonb,
  "type" text not null,
  "unit_amount" bigint,
  "billing_scheme" text not null,
  "created" bigint not null,
  "livemode" boolean not null,
  "lookup_key" text,
  "tiers_mode" text,
  "transform_quantity" jsonb,
  "unit_amount_decimal" double precision,

  "product" text not null references stripe.products on delete cascade,

  "updated_at" timestamptz not null
);