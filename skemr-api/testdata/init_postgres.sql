-- ./testdata/init_postgres.sql
-- Generic schema + seed data in public schema

SET search_path = public;

-- Customers
CREATE TABLE IF NOT EXISTS customers (
                                         id          BIGSERIAL PRIMARY KEY,
                                         email       TEXT        NOT NULL UNIQUE,
                                         name        TEXT        NOT NULL,
                                         created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Products
CREATE TABLE IF NOT EXISTS products (
                                        id          BIGSERIAL PRIMARY KEY,
                                        sku         TEXT        NOT NULL UNIQUE,
                                        name        TEXT        NOT NULL,
                                        price       NUMERIC(12,2) NOT NULL CHECK (price >= 0),
                                        active      BOOLEAN     NOT NULL DEFAULT TRUE,
                                        created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Orders
CREATE TABLE IF NOT EXISTS orders (
                                      id          BIGSERIAL PRIMARY KEY,
                                      customer_id BIGINT NOT NULL REFERENCES customers3(id) ON DELETE CASCADE,
                                      status      TEXT   NOT NULL DEFAULT 'pending'
                                          CHECK (status IN ('pending','paid','shipped','cancelled')),
                                      created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Order items (composite PK)
CREATE TABLE IF NOT EXISTS order_items (
                                           order_id    BIGINT NOT NULL REFERENCES orders(id)   ON DELETE CASCADE,
                                           product_id  BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
                                           quantity    INTEGER NOT NULL CHECK (quantity > 0),
                                           unit_price  NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
                                           PRIMARY KEY (order_id, product_id)
);

-- Helpful indexes
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);

-- -----------------------------
-- Seed data (idempotent)
-- -----------------------------

-- Customers
INSERT INTO customers3 (id, email, name)
VALUES
    (1, 'alice@example.com', 'Alice'),
    (2, 'bob@example.com',   'Bob')
ON CONFLICT (id) DO NOTHING;

-- Products
INSERT INTO products (id, sku, name, price)
VALUES
    (1, 'SKU-001', 'Widget',     19.99),
    (2, 'SKU-002', 'Gadget',     29.50),
    (3, 'SKU-003', 'Doohickey',   9.99)
ON CONFLICT (id) DO NOTHING;

-- Orders
INSERT INTO orders (id, customer_id, status)
VALUES
    (1, 1, 'paid'),
    (2, 2, 'pending')
ON CONFLICT (id) DO NOTHING;

-- Order items
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES
    (1, 1, 2, 19.99),
    (1, 3, 1,  9.99),
    (2, 2, 1, 29.50)
ON CONFLICT DO NOTHING;

-- -----------------------------
-- Reset sequences to max(id)
-- -----------------------------
SELECT setval(pg_get_serial_sequence('customers','id'),
              COALESCE((SELECT MAX(id) FROM customers3), 0), true);

SELECT setval(pg_get_serial_sequence('products','id'),
              COALESCE((SELECT MAX(id) FROM products), 0), true);

SELECT setval(pg_get_serial_sequence('orders','id'),
              COALESCE((SELECT MAX(id) FROM orders), 0), true);


-- Create analytics schema
CREATE SCHEMA IF NOT EXISTS analytics;

-- Create table in analytics schema
CREATE TABLE IF NOT EXISTS analyticss.page_views (
                                                    id         BIGSERIAL PRIMARY KEY,
                                                    customer_id BIGINT REFERENCES public.customers3(id) ON DELETE SET NULL,
                                                    page_url   TEXT NOT NULL,
                                                    viewed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed data for analytics.page_views
INSERT INTO analyticss.page_views (id, customer_id, page_url, viewed_at)
VALUES
    (1, 1, '/home',      now() - INTERVAL '2 days'),
    (2, 1, '/products',  now() - INTERVAL '1 day'),
    (3, 2, '/checkout',  now())
ON CONFLICT (id) DO NOTHING;

-- Reset sequence for analytics.page_views
SELECT setval(pg_get_serial_sequence('analytics.page_views','id'),
              COALESCE((SELECT MAX(id) FROM analyticss.page_views), 0), true);