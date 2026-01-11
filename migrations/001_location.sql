-- +goose Up
-- +goose StatementBegin

CREATE TABLE locations (
    id UUID PRIMARY KEY,
    name varchar(255) NOT NULL,
    latitude INTEGER NOT NULL,
    longitude INTEGER NOT NULL,
    rewardMultiplier DOUBLE PRECISION NOT NULL,
    type varchar(255) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd