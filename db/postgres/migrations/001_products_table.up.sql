/* https://stripe.com/docs/api/products/object */

create table if not exists "stripe"."products" (
  "id" text primary key,
  "object" text not null,
  "active" boolean not null,
  "description" text,
  "metadata" jsonb,
  "name" text not null,
  "created" bigint not null,
  "images" text[],
  "livemode" boolean not null,
  "package_dimensions" jsonb,
  "shippable" boolean not null,
  "statement_descriptor" text,
  "unit_label" text,
  "updated" bigint not null,
  "url" text,

  updated_at timestamptz not null
);

