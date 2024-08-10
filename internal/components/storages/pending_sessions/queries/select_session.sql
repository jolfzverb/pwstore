SELECT
    idempotency_token,
    session_id,
    nonce,
    state
FROM sessions_tmp
WHERE
    session_id=$1;
