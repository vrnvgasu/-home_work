-- +goose Up
-- +goose StatementBegin

UPDATE event
SET description = ''
WHERE description IS NULL;

UPDATE event
SET send_before = 0
WHERE send_before IS NULL;

ALTER TABLE event
    ALTER COLUMN description SET DEFAULT '',
    ALTER COLUMN description SET NOT NULL;
ALTER TABLE event
    ALTER COLUMN send_before SET DEFAULT 0,
    ALTER COLUMN send_before SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE event
    ALTER COLUMN description DROP NOT NULL;
ALTER TABLE event
    ALTER COLUMN send_before DROP NOT NULL;

-- +goose StatementEnd
