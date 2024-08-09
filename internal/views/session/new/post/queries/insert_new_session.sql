INSERT INTO sessions_tmp (
    idempotency_token
) VALUES (
    $1
) ON CONFLICT (idempotency_token) DO UPDATE SET idempotency_token = $1
RETURNING
   idempotency_token,
   session_id,
   nonce,
   state
