-- +goose Up
-- +goose StatementBegin
CREATE TABLE login_attempts (
    email VARCHAR(70) NOT NULL,
    remote_ip bytea NOT NULL,
    login_attempt_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (email, remote_ip, login_attempt_time)
);

CREATE INDEX idx_login_attempts_email ON login_attempts(email);
CREATE INDEX idx_login_attempts_remote_ip ON login_attempts(remote_ip);
CREATE INDEX idx_login_attempts_login_attempt_time ON login_attempts(login_attempt_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_login_attempts_email;
DROP INDEX IF EXISTS idx_login_attempts_remote_ip;
DROP INDEX IF EXISTS idx_login_attempts_login_attempt_time;

DROP TABLE IF EXISTS login_attempts;
-- +goose StatementEnd
