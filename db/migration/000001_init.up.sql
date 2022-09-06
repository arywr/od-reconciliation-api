-- Transaction model table
CREATE TABLE "transaction_models" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "transaction_id" SMALLINT,
    "merchant_transaction_id" SMALLINT,
    "channel_transaction_id" SMALLINT,
    "channel_code" SMALLINT,
    "channel_name" SMALLINT,
    "merchant_code" SMALLINT,
    "merchant_name" SMALLINT,
    "product_code" SMALLINT,
    "product_name" SMALLINT,
    "transaction_date" SMALLINT,
    "transaction_datetime" SMALLINT,
    "transaction_amount" SMALLINT,
    "settled_amount" SMALLINT,
    "transaction_date_format" VARCHAR,
    "transaction_datetime_format" VARCHAR,
    "row_start_at" SMALLINT DEFAULT 0,
    "file_type" VARCHAR DEFAULT 'XLSX' -- CSV or XLSX
);

INSERT INTO transaction_models (name, transaction_id, merchant_transaction_id, channel_transaction_id, channel_code, channel_name, merchant_code, merchant_name, product_code, product_name, transaction_date, transaction_datetime, transaction_amount, settled_amount, transaction_date_format, transaction_datetime_format, row_start_at, file_type)
VALUES ('Data Model For T-Money', 2, NULL, NULL, NULL, NULL, 11, 12, 13, 14, 30, 30, 17, 17, '2006-01-02 15:04:05', '2006-01-02 15:04:05', 7, 'XLSX');

INSERT INTO transaction_models (name, transaction_id, merchant_transaction_id, channel_transaction_id, channel_code, channel_name, merchant_code, merchant_name, product_code, product_name, transaction_date, transaction_datetime, transaction_amount, settled_amount, transaction_date_format, transaction_datetime_format, row_start_at, file_type)
VALUES ('Data Model For ALTO (T-Money)', 1, 1, NULL, 1, 1, 1, 1, 1, 1, 8, 8, 19, 19, '02/01/2006 15:04:05', '2006-01-02 15:04:05', 1, 'XLSX');

-- Product tables
CREATE TABLE "products" (
    "id" BIGSERIAL PRIMARY KEY,
    "transaction_model_id" BIGINT DEFAULT NULL,
    "product_name" VARCHAR NOT NULL,
    "product_code" VARCHAR NOT NULL,
    "product_has_sub" BOOLEAN NOT NULL
);

CREATE TABLE "sub_products" (
    "id" BIGSERIAL PRIMARY KEY,
    "product_id" BIGINT NOT NULl,
    "sub_product_name" VARCHAR NOT NULL,
    "sub_product_code" VARCHAR NOT NULL
);

-- Inserting data to products table
INSERT INTO products (product_name, transaction_model_id, product_code, product_has_sub) VALUES ('T-Money', 1, '001', FALSE);
INSERT INTO products (product_name, transaction_model_id, product_code, product_has_sub) VALUES ('QREN', 1, '002', TRUE);
INSERT INTO products (product_name, transaction_model_id, product_code, product_has_sub) VALUES ('MPS', 1, '003', FALSE);

INSERT INTO sub_products (product_id, sub_product_name, sub_product_code) VALUES (2, 'INTERACTIVE', '001001');
INSERT INTO sub_products (product_id, sub_product_name, sub_product_code) VALUES (2, 'METRANET', '001002');
INSERT INTO sub_products (product_id, sub_product_name, sub_product_code) VALUES (2, 'TLT PARKING', '001003');

CREATE TABLE "merchants" (
    "id" BIGSERIAL PRIMARY KEY,
    "product_id" BIGINT NOT NULL,
    "transaction_model_id" INT DEFAULT NULL,
    "merchant_name" VARCHAR NOT NULL,
    "merchant_code" VARCHAR NOT NULL,
    "merchant_has_sub" BOOLEAN NOT NULL
);

CREATE TABLE "sub_merchants" (
    "id" BIGSERIAL PRIMARY KEY,
    "merchant_id" BIGINT NOT NULL,
    "sub_merchant_name" VARCHAR NOT NULL,
    "sub_merchant_code" VARCHAR NOT NULL
);

-- Inserting data to merchant table
INSERT INTO merchants (product_id, merchant_name, merchant_code, merchant_has_sub) VALUES (2, 'JALIN', '001', FALSE);
INSERT INTO merchants (product_id, merchant_name, merchant_code, merchant_has_sub) VALUES (1, 'ALTO', '002', FALSE);
INSERT INTO merchants (product_id, merchant_name, merchant_code, merchant_has_sub) VALUES (1, 'Metranet-Espay', '003', FALSE);
INSERT INTO merchants (product_id, merchant_name, merchant_code, merchant_has_sub) VALUES (1, 'Metranet', '004', FALSE);

-- Reconcile rule table
-- Column conditions            = EQUAL, CONTAINS, NOT EQUAL
-- Rule mandatory conditions    = REQUIRED, OPTIONAL
CREATE TABLE "reconcile_rules" (
    "id" BIGSERIAL PRIMARY KEY,
    "product_id" BIGINT NOT NULL,
    "platform_id" BIGINT NOT NULL,
    "product_column_field" VARCHAR DEFAULT NULL,
    "product_column_conditions" VARCHAR DEFAULT NULL,
    "product_column_value" VARCHAR DEFAULT NULL,
    "platform_column_field" VARCHAR DEFAULT NULL,
    "platform_column_conditions" VARCHAR DEFAULT NULL,
    "platform_column_value" VARCHAR DEFAULT NULL,
    "rule_mandatory" VARCHAR DEFAULT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

INSERT INTO reconcile_rules (product_id, platform_id, product_column_field, product_column_conditions, product_column_value, platform_column_field, platform_column_conditions, platform_column_value, rule_mandatory) 
VALUES (1, 4, 'merchant_name', 'EQUAL', 'METRANET', NULL, NULL, NULL, 'REQUIRED');

INSERT INTO reconcile_rules (product_id, platform_id, product_column_field, product_column_conditions, product_column_value, platform_column_field, platform_column_conditions, platform_column_value, rule_mandatory) 
VALUES (2, 1, 'transaction_key', 'EQUAL', 'ANY', 'transaction_key', 'EQUAL', 'ANY', 'REQUIRED');

INSERT INTO reconcile_rules (product_id, platform_id, product_column_field, product_column_conditions, product_column_value, platform_column_field, platform_column_conditions, platform_column_value, rule_mandatory) 
VALUES (1, 3, 'merchant_name', 'EQUAL', 'ESPAY_TMONEY', NULL, NULL, NULL, 'REQUIRED');

INSERT INTO reconcile_rules (product_id, platform_id, product_column_field, product_column_conditions, product_column_value, platform_column_field, platform_column_conditions, platform_column_value, rule_mandatory) 
VALUES (1, 2, NULL, NULL, NULL, 'merchant_name', 'EQUAL', 'ANY', 'REQUIRED');

CREATE TABLE "transaction_statuses" (
    "id" BIGSERIAL PRIMARY KEY,
    "status_name" VARCHAR UNIQUE NOT NULL,
    "status_description" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

INSERT INTO transaction_statuses (status_name, status_description) VALUES ('UNKNOWN', 'This is description');
INSERT INTO transaction_statuses (status_name, status_description) VALUES ('PARTIALLY-MATCH', 'This is description');
INSERT INTO transaction_statuses (status_name, status_description) VALUES ('FULLY-MATCH', 'This is description');
INSERT INTO transaction_statuses (status_name, status_description) VALUES ('DISPUTE', 'This is description');

CREATE TABLE "transaction_types" (
    "id" BIGSERIAL PRIMARY KEY,
    "type_name" VARCHAR UNIQUE NOT NULL,
    "type_description" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

INSERT INTO transaction_types (type_name, type_description) VALUES ('CASH-OUT', 'This is description');
INSERT INTO transaction_types (type_name, type_description) VALUES ('CASH-IN', 'This is description');

CREATE TABLE "progress_event_types" (
    "id" BIGSERIAL PRIMARY KEY,
    "progress_event_type_name" VARCHAR UNIQUE NOT NULL,
    "progress_event_type_description" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

INSERT INTO progress_event_types (progress_event_type_name, progress_event_type_description) VALUES ('READ', 'This is description');
INSERT INTO progress_event_types (progress_event_type_name, progress_event_type_description) VALUES ('RECONCILE', 'This is description');

CREATE TABLE "progress_events" (
    "id" BIGSERIAL PRIMARY KEY,
    "progress_event_type_id" SMALLINT NOT NULL,
    "progress_name" VARCHAR NOT NULL,
    "status" VARCHAR NOT NULL,
    "percentage" FLOAT DEFAULT 0 NOT NULL,
    "file" VARCHAR NOT NULL DEFAULT 'Unknown File',
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "product_transactions" (
    "id" BIGSERIAL PRIMARY KEY,
    "product_id" INT NOT NULL,
    "sub_product_id" INT DEFAULT NULL,
    "platform_id" INT DEFAULT NULL,
    "sub_platform_id" INT DEFAULT NULL,
    "transaction_status_id" SMALLINT NOT NULL,
    "transaction_type_id" SMALLINT NOT NULL,
    "progress_event_id" INT DEFAULT NULL,
    "product_transaction_id" VARCHAR DEFAULT NULL,
    "merchant_transaction_id" VARCHAR DEFAULT NULL,
    "transaction_id" VARCHAR NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" TIMESTAMP NOT NULL DEFAULT(now()),
    "channel_code" VARCHAR,
    "channel_name" VARCHAR,
    "merchant_code" VARCHAR,
    "merchant_name" VARCHAR,
    "product_code" VARCHAR,
    "product_name" VARCHAR,
    "collected_amount" FLOAT DEFAULT 0 NOT NULL,
    "settled_amount" FLOAT DEFAULT 0 NOT NULL,
    "reconcile_at" TIMESTAMP DEFAULT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "merchant_transactions" (
    "id" BIGSERIAL PRIMARY KEY,
    "product_id" INT NOT NULL,
    "sub_product_id" INT DEFAULT NULL,
    "platform_id" INT DEFAULT NULL,
    "sub_platform_id" INT DEFAULT NULL,
    "transaction_status_id" SMALLINT NOT NULL,
    "transaction_type_id" SMALLINT NOT NULL,
    "progress_event_id" INT DEFAULT NULL,
    "product_transaction_id" VARCHAR DEFAULT NULL,
    "merchant_transaction_id" VARCHAR DEFAULT NULL,
    "transaction_id" VARCHAR NOT NULL,
    "transaction_date" date NOT NULL DEFAULT(now()),
    "transaction_datetime" TIMESTAMP NOT NULL DEFAULT(now()),
    "channel_code" VARCHAR,
    "channel_name" VARCHAR,
    "merchant_code" VARCHAR,
    "merchant_name" VARCHAR,
    "product_code" VARCHAR,
    "product_name" VARCHAR,
    "collected_amount" FLOAT DEFAULT 0 NOT NULL,
    "settled_amount" FLOAT DEFAULT 0 NOT NULL,
    "reconcile_at" TIMESTAMP DEFAULT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT(now()),
    "deleted_at" TIMESTAMP DEFAULT NULL
);

ALTER TABLE "product_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "transaction_statuses" ("id") ON DELETE CASCADE;
ALTER TABLE "product_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "transaction_types" ("id")  ON DELETE CASCADE;
ALTER TABLE "product_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "progress_events" ("id");

ALTER TABLE "merchant_transactions" ADD FOREIGN KEY ("transaction_status_id") REFERENCES "transaction_statuses" ("id") ON DELETE CASCADE;
ALTER TABLE "merchant_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "transaction_types" ("id")  ON DELETE CASCADE;
ALTER TABLE "merchant_transactions" ADD FOREIGN KEY ("progress_event_id") REFERENCES "progress_events" ("id")  ON DELETE CASCADE;

ALTER TABLE "progress_events" ADD FOREIGN KEY ("progress_event_type_id") REFERENCES "progress_events" ("id")  ON DELETE CASCADE;