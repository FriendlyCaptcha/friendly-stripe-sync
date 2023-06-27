create table if not exists "stripe"."sync_state" (
  "id" text primary key,
  "last_event" bigint not null
);