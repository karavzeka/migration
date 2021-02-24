CREATE TABLE IF NOT EXISTS section(
   section_id serial PRIMARY KEY,
   title VARCHAR (50) UNIQUE NOT NULL
);
