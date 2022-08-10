CREATE TYPE progress_status_enum AS ENUM('on process', 'completed', 'failed');

CREATE TABLE "od_transaction_statuses" (
    "id" bigserial PRIMARY KEY,
    "status_name" varchar NOT NULL,
    "status_description" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

CREATE TABLE "od_transaction_types" (
    "id" bigserial PRIMARY KEY,
    "type_name" varchar NOT NULL,
    "type_description" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

CREATE TABLE "od_progress_event_types" (
    "id" bigserial PRIMARY KEY,
    "progress_event_type_name" varchar NOT NULL,
    "progress_event_type_description" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

CREATE TABLE "od_progress_events" (
    "id" bigserial PRIMARY KEY,
    "progress_event_type_id" smallint NOT NULL,
    "progress_name" varchar NOT NULL,
    "status" progress_status_enum DEFAULT NULL,
    "percentage" float DEFAULT 0 NOT NULL,
    "file" varchar,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

CREATE TABLE "od_product_transactions" (
    "id" bigserial PRIMARY KEY,
    "transaction_status_id" smallint NOT NULL,
    "transaction_type_id" smallint NOT NULL,
    "progress_event_id" smallint NOT NULL,
    "product_transaction_id" varchar DEFAULT NULL,
    "merchant_transaction_id" varchar DEFAULT NULL,
    "channel_transaction_id" varchar DEFAULT NULL,
    "owner_id" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" timestamptz NOT NULL DEFAULT(now()),
    "collected_amount" float DEFAULT 0,
    "settled_amount" float DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

CREATE TABLE "od_merchant_transactions" (
    "id" bigserial PRIMARY KEY,
    "transaction_status_id" smallint NOT NULL,
    "transaction_type_id" smallint NOT NULL,
    "progress_event_id" smallint NOT NULL,
    "merchant_transaction_id" varchar DEFAULT NULL,
    "owner_id" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" timestamptz NOT NULL DEFAULT(now()),
    "collected_amount" float DEFAULT 0,
    "settled_amount" float DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

CREATE TABLE "od_channel_transactions" (
    "id" bigserial PRIMARY KEY,
    "transaction_status_id" smallint NOT NULL,
    "transaction_type_id" smallint NOT NULL,
    "progress_event_id" smallint NOT NULL,
    "channel_transaction_id" varchar DEFAULT NULL,
    "owner_id" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" timestamptz NOT NULL DEFAULT(now()),
    "collected_amount" float DEFAULT 0,
    "settled_amount" float DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "updated_at" timestamptz NOT NULL DEFAULT(now()),
    "deleted_at" timestamptz NOT NULL DEFAULT NULL
);

ALTER TABLE "od_product_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "od_transaction_statuses" ("id");
ALTER TABLE "od_product_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "od_transaction_types" ("id");
ALTER TABLE "od_product_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "od_progress_events" ("id");

ALTER TABLE "od_merchant_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "od_transaction_statuses" ("id");
ALTER TABLE "od_merchant_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "od_transaction_types" ("id");
ALTER TABLE "od_merchant_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "od_progress_events" ("id");

ALTER TABLE "od_channel_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "od_transaction_statuses" ("id");
ALTER TABLE "od_channel_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "od_transaction_types" ("id");
ALTER TABLE "od_channel_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "od_progress_events" ("id");

ALTER TABLE "od_progress_events" ADD FOREIGN KEY ("progress_event_type_id") REFERENCES "od_progress_events" ("id");