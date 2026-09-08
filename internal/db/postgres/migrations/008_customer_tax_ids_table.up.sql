/* https://stripe.com/docs/api/customer_tax_ids/object */

create table if not exists "stripe"."customer_tax_ids" (
  "id" text primary key,
  "object" text not null,
  "country" text,
  "created" bigint not null,
  "livemode" boolean not null,
  "type" text not null,
  "value" text not null,
  "verification" jsonb,
  "customer" text not null references "stripe"."customers" on delete cascade,

  "updated_at" timestamptz not null
);

create index if not exists "customer_tax_ids_customer_idx" on "stripe"."customer_tax_ids" ("customer");
