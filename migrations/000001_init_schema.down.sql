DROP INDEX IF EXISTS idx_products_reception_id_datetime;
DROP INDEX IF EXISTS idx_products_reception_id;
DROP TABLE IF EXISTS products;

DROP INDEX IF EXISTS idx_receptions_pvz_id_status_open;
DROP INDEX IF EXISTS idx_receptions_date_time;
DROP INDEX IF EXISTS idx_receptions_status;
DROP INDEX IF EXISTS idx_receptions_pvz_id;
DROP TABLE IF EXISTS receptions;

DROP INDEX IF EXISTS idx_pvz_city;
DROP TABLE IF EXISTS pvz;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS product_type;
DROP TYPE IF EXISTS reception_status;
DROP TYPE IF EXISTS pvz_city;
DROP TYPE IF EXISTS user_role;