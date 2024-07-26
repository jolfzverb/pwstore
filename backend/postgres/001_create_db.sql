CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE items (
    id TEXT PRIMARY KEY NOT NULL default uuid_generate_v4(),
    name TEXT NOT NULL,
    price FLOAT NOT NULL
);