CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE sessions (
    session_id TEXT PRIMARY KEY NOT NULL default uuid_generate_v4(),
    token TEXT PRIMARY KEY NOT NULL default uuid_generate_v4()
);
CREATE UNIQUE INDEX sessions_token on sessions(token);

CREATE TABLE sessions_tmp (
   idempotency_token TEXT PRIMARY KEY NOT NULL,
   session_id TEXT NOT NULL default uuid_generate_v4(),
   nonce TEXT NOT NULL default uuid_generate_v4(),
   state TEXT NOT NULL default uuid_generate_v4()
);
CREATE UNIQUE INDEX sessions_tmp_session_id on sessions_tmp(session_id);
