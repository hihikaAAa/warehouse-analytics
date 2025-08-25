-- +goose Up
CREATE TABLE IF NOT EXISTS mp_platform (
  id serial PRIMARY KEY,
  code text UNIQUE NOT NULL,
  name text NOT NULL
);

CREATE TABLE IF NOT EXISTS product (
  id serial PRIMARY KEY,
  sku text UNIQUE NOT NULL,
  title text,
  barcode text,
  category text
);

CREATE TABLE IF NOT EXISTS listing (
  id serial PRIMARY KEY,
  product_id int REFERENCES product(id),
  platform_id int REFERENCES mp_platform(id),
  mp_sku text,
  UNIQUE(product_id, platform_id),
  UNIQUE(platform_id, mp_sku)
);

CREATE TABLE IF NOT EXISTS orders (
  id bigserial PRIMARY KEY,
  platform_id int REFERENCES mp_platform(id),
  mp_order_id text NOT NULL,
  created_at timestamptz NOT NULL,
  status text,
  buyer_region text,
  UNIQUE(platform_id, mp_order_id)
);
CREATE INDEX IF NOT EXISTS idx_orders_platform_created ON orders(platform_id, created_at);

CREATE TABLE IF NOT EXISTS order_items (
  id bigserial PRIMARY KEY,
  order_id bigint REFERENCES orders(id) ON DELETE CASCADE,
  listing_id int REFERENCES listing(id),
  qty int NOT NULL,
  price numeric(12,2) NOT NULL,
  discount numeric(12,2) NOT NULL DEFAULT 0,
  revenue numeric(12,2) GENERATED ALWAYS AS (qty*(price-discount)) STORED,
  UNIQUE(order_id, listing_id)
);

CREATE TABLE IF NOT EXISTS shipments (
  id bigserial PRIMARY KEY,
  platform_id int REFERENCES mp_platform(id),
  mp_shipment_id text NOT NULL,
  status text,
  shipped_at timestamptz,
  delivered_at timestamptz
);

CREATE TABLE IF NOT EXISTS returns (
  id bigserial PRIMARY KEY,
  order_item_id bigint REFERENCES order_items(id) ON DELETE CASCADE,
  reason text,
  returned_at timestamptz
);

CREATE TABLE IF NOT EXISTS stock (
  id bigserial PRIMARY KEY,
  listing_id int REFERENCES listing(id),
  warehouse text,
  qty int NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE(listing_id, warehouse)
);
CREATE INDEX IF NOT EXISTS idx_stock_listing_updated ON stock(listing_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS stock_movements (
  id bigserial PRIMARY KEY,
  listing_id int REFERENCES listing(id),
  warehouse text,
  delta int NOT NULL,
  reason text,
  occurred_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS kpi_sales_daily (
  day date NOT NULL,
  listing_id int REFERENCES listing(id),
  orders int,
  units int,
  revenue numeric(12,2),
  returns int,
  PRIMARY KEY (day, listing_id)
);

CREATE TABLE IF NOT EXISTS etl_cursors (
  id serial PRIMARY KEY,
  platform_id int REFERENCES mp_platform(id),
  resource text NOT NULL,
  cursor_value text NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE(platform_id, resource)
);

-- +goose Down
DROP TABLE IF EXISTS etl_cursors;
DROP TABLE IF EXISTS kpi_sales_daily;
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS stock;
DROP TABLE IF EXISTS returns;
DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS listing;
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS mp_platform;
