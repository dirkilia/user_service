CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	surname TEXT NOT NULL,
	patronymic TEXT,
	age INTEGER,
	gender TEXT,
	nationality TEXT,
	CHECK (age > 0)
);