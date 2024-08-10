INSERT INTO sessions (
    session_id,
    subject,
    email,
    id_token
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (session_id) DO UPDATE SET idempotency_token = $1
RETURNING
    session_id,
    subject,
    email,
    id_token,
    token
