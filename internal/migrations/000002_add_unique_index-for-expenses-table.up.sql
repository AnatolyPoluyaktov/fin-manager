ALTER TABLE expenses
ADD CONSTRAINT unique_category_date_amount UNIQUE (category_id, action_date, amount);