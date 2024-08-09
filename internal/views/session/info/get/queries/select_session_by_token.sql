SELECT
    session_id
FROM sessions
WHERE
    token=$1;
