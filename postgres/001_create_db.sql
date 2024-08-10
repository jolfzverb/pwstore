CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE sessions (
    session_id TEXT PRIMARY KEY NOT NULL,
    subject TEXT NOT NULL,
    email TEXT NOT NULL,
    id_token TEXT NOT NULL,
    token TEXT NOT NULL DEFAULT UUID_GENERATE_V4()
);
CREATE UNIQUE INDEX sessions_token ON sessions(token);
CREATE INDEX sessions_subject ON sessions(subject);
CREATE INDEX sessions_email ON sessions(email);

CREATE TABLE pending_sessions (
   idempotency_token TEXT PRIMARY KEY NOT NULL,
   session_id TEXT NOT NULL DEFAULT UUID_GENERATE_V4(),
   nonce TEXT NOT NULL DEFAULT UUID_GENERATE_V4(),
   state TEXT NOT NULL DEFAULT UUID_GENERATE_V4()
);
CREATE UNIQUE INDEX sessions_tmp_session_id ON pending_sessions(session_id);
