create table if not exists "stripe"."coupons" (
  "id" text primary key,
  "object" text not null,
  "amount_off" bigint,
  "created" bigint not null,
  "currency" text,
  "duration" text not null,
  "duration_in_months" bigint,
  "livemode" boolean not null,
  "max_redemptions" bigint,
  "metadata" jsonb,
  "name" text not null,
  "percent_off" float8,
  "redeem_by" bigint,
  "times_redeemed" bigint,
  "valid" boolean not null,
  "applies_to" jsonb,
  "updated_at" timestamptz not null
);