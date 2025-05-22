-- Файл: migrations/000001_init_schema.up.sql

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('employee', 'moderator');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pvz_city') THEN
        CREATE TYPE pvz_city AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reception_status') THEN
        CREATE TYPE reception_status AS ENUM ('in_progress', 'close'); -- В ТЗ close, у вас closed. Привожу к ТЗ.
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type') THEN
        CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS users ( -- Было: users( , EXIST
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Убрал UNIQUE, пароли не должны быть уникальными
    role user_role NOT NULL
);

CREATE TABLE IF NOT EXISTS pvz ( -- Было: EXIST
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- В вашей миграции registration_data, в ТЗ registration_date
    city pvz_city NOT NULL
);

CREATE TABLE IF NOT EXISTS receptions ( -- Было: EXIST
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- В вашей миграции data_time, в ТЗ date_time
    pvz_id UUID NOT NULL,
    status reception_status NOT NULL DEFAULT 'in_progress', -- Было "in_progress"
    FOREIGN KEY (pvz_id) REFERENCES pvz(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS products ( -- Было: product, EXIST, нет ;
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- В вашей миграции data_time, в ТЗ date_time
    type product_type NOT NULL,
    reception_id UUID NOT NULL, -- В вашей миграции receptions_id, в ТЗ reception_id
    FOREIGN KEY (reception_id) REFERENCES receptions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Это внутренний ID записи
    jti VARCHAR(255) NOT NULL UNIQUE, -- JTI должен быть уникальным, чтобы предотвратить replay атаки на refresh token
    user_id VARCHAR(255) NOT NULL, -- Изменил на VARCHAR для совместимости с uuid.New().String() напрямую, либо используйте UUID и конвертируйте
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    ip VARCHAR(45),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);