-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE payment_method AS ENUM (
    'UNKNOWN',
    'CARD',
    'SBP',
    'CREDIT_CARD',
    'INVESTOR_MONEY'
);

CREATE TYPE order_status AS ENUM (
    'UNKNOWN',
    'PENDING_PAYMENT',
    'PAID',
    'CANCELLED'
);

create table orders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    parts_ids uuid[] NOT NULL DEFAULT '{}',
    total_price REAL NOT NULL DEFAULT 0,
    transaction_id uuid,
    payment_method payment_method NOT NULL DEFAULT 'UNKNOWN',
    status order_status NOT NULL DEFAULT 'UNKNOWN'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS payment_method;
-- +goose StatementEnd