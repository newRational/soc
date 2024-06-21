-- +goose Up
-- +goose StatementBegin
create table orders(
    id          BIGSERIAL PRIMARY KEY NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    level       BIGINT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table orders;
-- +goose StatementEnd