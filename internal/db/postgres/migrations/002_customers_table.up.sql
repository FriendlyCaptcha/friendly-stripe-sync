/* https://stripe.com/docs/api/products/object */

create table if not exists "stripe"."customers" (
  "id" text primary key,
  "object" text not null,
  "address" jsonb,
  "description" text,
  "email" text,
  "metadata" jsonb,
  "name" text,
  "phone" text,
  "shipping" jsonb,
  "balance" bigint,
  "created" bigint not null,
  "currency" text,
  "default_source" text,
  "delinquent" boolean not null,
  "discount" jsonb,
  "invoice_prefix" text not null,
  "invoice_settings" jsonb,
  "livemode" boolean not null,
  "next_invoice_sequence" bigint not null,
  "preferred_locales" text[],
  "tax_exempt" text not null,
  "deleted" boolean not null,

  "updated_at" timestamptz not null
);

