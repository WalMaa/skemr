-- ./testdata/init_postgres.sql
-- Fully schema-qualified schema + seed data


--- Create public schema if not exists ---
CREATE SCHEMA IF NOT EXISTS public;
-- Public tables
CREATE TABLE IF NOT EXISTS public.customers
(
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT        NOT NULL UNIQUE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT pg_catalog.now()
);

CREATE TABLE IF NOT EXISTS public.products
(
    id         BIGSERIAL PRIMARY KEY,
    sku        TEXT           NOT NULL UNIQUE,
    name       TEXT           NOT NULL,
    price      NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    active     BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT pg_catalog.now()
);

CREATE TABLE IF NOT EXISTS public.orders
(
    id          BIGSERIAL PRIMARY KEY,
    customer_id BIGINT      NOT NULL
        REFERENCES public.customers (id) ON DELETE CASCADE,
    status      TEXT        NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'paid', 'shipped', 'cancelled')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT pg_catalog.now()
);

CREATE TABLE IF NOT EXISTS public.order_items
(
    order_id   BIGINT         NOT NULL
        REFERENCES public.orders (id) ON DELETE CASCADE,
    product_id BIGINT         NOT NULL
        REFERENCES public.products (id) ON DELETE RESTRICT,
    quantity   INTEGER        NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(12, 2) NOT NULL CHECK (unit_price >= 0),
    PRIMARY KEY (order_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id
    ON public.orders (customer_id);

-- Seed data for public.customers
INSERT INTO public.customers (id, email, name)
VALUES (1, 'alice@example.com', 'Alice'),
       (2, 'bob@example.com', 'Bob')
ON CONFLICT (id) DO NOTHING;

-- Seed data for public.products
INSERT INTO public.products (id, sku, name, price)
VALUES (1, 'SKU-001', 'Widget', 19.99),
       (2, 'SKU-002', 'Gadget', 29.50),
       (3, 'SKU-003', 'Doohickey', 9.99)
ON CONFLICT (id) DO NOTHING;

-- Seed data for public.orders
INSERT INTO public.orders (id, customer_id, status)
VALUES (1, 1, 'paid'),
       (2, 2, 'pending')
ON CONFLICT (id) DO NOTHING;

-- Seed data for public.order_items
INSERT INTO public.order_items (order_id, product_id, quantity, unit_price)
VALUES (1, 1, 2, 19.99),
       (1, 3, 1, 9.99),
       (2, 2, 1, 29.50)
ON CONFLICT DO NOTHING;

-- Reset sequences for public tables
SELECT pg_catalog.setval(
               pg_catalog.pg_get_serial_sequence('public.customers', 'id'),
               COALESCE((SELECT MAX(id) FROM public.customers), 0),
               true
       );

SELECT pg_catalog.setval(
               pg_catalog.pg_get_serial_sequence('public.products', 'id'),
               COALESCE((SELECT MAX(id) FROM public.products), 0),
               true
       );

SELECT pg_catalog.setval(
               pg_catalog.pg_get_serial_sequence('public.orders', 'id'),
               COALESCE((SELECT MAX(id) FROM public.orders), 0),
               true
       );

-- Analytics schema
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.page_views
(
    id          BIGSERIAL PRIMARY KEY,
    customer_id BIGINT
                            REFERENCES public.customers (id) ON DELETE SET NULL,
    page_url    TEXT        NOT NULL,
    viewed_at   TIMESTAMPTZ NOT NULL DEFAULT pg_catalog.now()
);

-- Seed data for analytics.page_views
INSERT INTO analytics.page_views (id, customer_id, page_url, viewed_at)
VALUES (1, 1, '/home', pg_catalog.now() - INTERVAL '2 days'),
       (2, 1, '/products', pg_catalog.now() - INTERVAL '1 day'),
       (3, 2, '/checkout', pg_catalog.now())
ON CONFLICT (id) DO NOTHING;

-- Reset sequence for analytics.page_views
SELECT pg_catalog.setval(
               pg_catalog.pg_get_serial_sequence('analytics.page_views', 'id'),
               COALESCE((SELECT MAX(id) FROM analytics.page_views), 0),
               true
       );