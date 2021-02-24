CREATE TABLE IF NOT EXISTS product(
   product_id serial PRIMARY KEY,
   title VARCHAR (50) UNIQUE NOT NULL
);
