INSERT INTO sessions (
    session_id
) VALUES (
    $1
) RETURNING
    session_id,
    token
