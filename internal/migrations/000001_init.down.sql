-- Обратная миграция: удаление таблицы и типов ENUM
DROP TABLE IF EXISTS expenses;
DROP TYPE IF EXISTS currency_enum;
DROP TABLE IF EXISTS categories;