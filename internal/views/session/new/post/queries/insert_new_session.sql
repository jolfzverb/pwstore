INSERT INTO sessions_tmp (
    idempotency_token
) VALUES (
    $1
) RETURNING
    session_id
