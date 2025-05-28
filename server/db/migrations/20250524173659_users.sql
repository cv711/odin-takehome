-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email varchar(70) NOT NULL,
  password_hash varchar(150) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_users_updated_unix_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION updated_at_unix_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_users_updated_unix_timestamp ON users;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
