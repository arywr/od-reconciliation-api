CREATE TABLE "transaction_statuses" (
    "id" bigserial PRIMARY KEY,
    "status_name" varchar UNIQUE NOT NULL,
    "status_description" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

CREATE TABLE "transaction_types" (
    "id" bigserial PRIMARY KEY,
    "type_name" varchar UNIQUE NOT NULL,
    "type_description" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

CREATE TABLE "progress_event_types" (
    "id" bigserial PRIMARY KEY,
    "progress_event_type_name" varchar UNIQUE NOT NULL,
    "progress_event_type_description" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

CREATE TABLE "progress_events" (
    "id" bigserial PRIMARY KEY,
    "progress_event_type_id" smallint NOT NULL,
    "progress_name" varchar NOT NULL,
    "status" varchar NOT NULL,
    "percentage" float DEFAULT 0 NOT NULL,
    "file" varchar NOT NULL DEFAULT 'Unknown File',
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

CREATE TABLE "product_transactions" (
    "id" bigserial PRIMARY KEY,
    "transaction_status_id" smallint NOT NULL,
    "transaction_type_id" smallint NOT NULL,
    "progress_event_id" int DEFAULT NULL,
    "product_transaction_id" varchar DEFAULT NULL,
    "merchant_transaction_id" varchar DEFAULT NULL,
    "channel_transaction_id" varchar DEFAULT NULL,
    "owner_id" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" timestamp NOT NULL DEFAULT(now()),
    "collected_amount" float DEFAULT 0 NOT NULL,
    "settled_amount" float DEFAULT 0 NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

CREATE TABLE "merchant_transactions" (
    "id" bigserial PRIMARY KEY,
    "transaction_status_id" smallint NOT NULL,
    "transaction_type_id" smallint NOT NULL,
    "progress_event_id" smallint NOT NULL,
    "merchant_transaction_id" varchar DEFAULT NULL,
    "owner_id" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" timestamp NOT NULL DEFAULT(now()),
    "collected_amount" float DEFAULT 0 NOT NULL,
    "settled_amount" float DEFAULT 0 NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

CREATE TABLE "channel_transactions" (
    "id" bigserial PRIMARY KEY,
    "transaction_status_id" smallint NOT NULL,
    "transaction_type_id" smallint NOT NULL,
    "progress_event_id" smallint NOT NULL,
    "channel_transaction_id" varchar DEFAULT NULL,
    "owner_id" varchar NOT NULL,
    "transaction_id" varchar NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" timestamp NOT NULL DEFAULT(now()),
    "collected_amount" float DEFAULT 0 NOT NULL,
    "settled_amount" float DEFAULT 0 NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT(now()),
    "updated_at" timestamp NOT NULL DEFAULT(now()),
    "deleted_at" timestamp DEFAULT NULL
);

ALTER TABLE "product_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "transaction_statuses" ("id") ON DELETE CASCADE;
ALTER TABLE "product_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "transaction_types" ("id")  ON DELETE CASCADE;
ALTER TABLE "product_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "progress_events" ("id");

ALTER TABLE "merchant_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "transaction_statuses" ("id")  ON DELETE CASCADE;
ALTER TABLE "merchant_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "transaction_types" ("id")  ON DELETE CASCADE;
ALTER TABLE "merchant_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "progress_events" ("id")  ON DELETE CASCADE;

ALTER TABLE "channel_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "transaction_statuses" ("id")  ON DELETE CASCADE;
ALTER TABLE "channel_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "transaction_types" ("id")  ON DELETE CASCADE;
ALTER TABLE "channel_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "progress_events" ("id")  ON DELETE CASCADE;

ALTER TABLE "progress_events" ADD FOREIGN KEY ("progress_event_type_id") REFERENCES "progress_events" ("id")  ON DELETE CASCADE;