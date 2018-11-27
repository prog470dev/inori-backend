
-- +migrate Up
ALTER TABLE reservations ADD departure_time DATETIME NOT NULL;

-- +migrate Down
ALTER TABLE reservations DROP COLUMN departure_time;
