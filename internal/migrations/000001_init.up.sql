CREATE TYPE currency_enum AS ENUM ('USD', 'EUR', 'RUB', 'GBP', 'JPY');
-- USD - Доллар США
-- EUR - Евро
-- RUB - Российский рубль
-- GBP - Британский фунт
-- JPY - Японская иена

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    currency currency_enum NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    action_date TIMESTAMP NOT NULL DEFAULT NOW(),
    note TEXT
);
