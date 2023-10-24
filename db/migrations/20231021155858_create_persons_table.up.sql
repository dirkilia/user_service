CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY NOT NULL,
	name TEXT,
	surname TEXT,
	patronymic TEXT,
	age INTEGER,
	gender TEXT,
	nationality TEXT,
	CHECK (age > 0)
);