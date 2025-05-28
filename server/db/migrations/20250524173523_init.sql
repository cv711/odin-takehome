-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION updated_at_unix_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := EXTRACT(EPOCH FROM NOW())::BIGINT;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION updated_at_unix_timestamp;
-- +goose StatementEnd
