SELECT
    idempotency_token,
    session_id,
    nonce,
    state
FROM pending_sessions
WHERE
    session_id=$1;
