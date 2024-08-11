SELECT
    session_id,
    subject,
    email,
    id_token,
    token
FROM sessions
WHERE
    session_id=$1 AND token=$2;
