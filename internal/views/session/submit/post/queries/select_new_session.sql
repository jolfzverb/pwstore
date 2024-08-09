SELECT
    session_id
FROM sessions_tmp
WHERE
    session_id=$1;
