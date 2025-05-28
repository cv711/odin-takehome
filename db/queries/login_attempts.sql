-- name: LogAttempt :one
INSERT INTO login_attempts (email, remote_ip) VALUES ($1, $2) RETURNING *;

-- name: GetCounts :one
SELECT
    COUNT(*) as global_count,
    COALESCE(SUM(CASE WHEN email = $1 THEN 1 ELSE 0 END),0)::BIGINT as email_count,
    COALESCE(SUM(CASE WHEN remote_ip = $2 THEN 1 ELSE 0 END),0)::BIGINT as ip_count
FROM login_attempts
WHERE login_attempt_time >= NOW() - INTERVAL '10 SECOND';

