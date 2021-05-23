-- +goose Up
-- +goose StatementBegin
CREATE TABLE "events" (
    "id" SERIAL NOT NULL,
    "title" VARCHAR(128) NOT NULL,
    "client_id" INT NOT NULL,
    "year" SMALLINT NOT NULL,
    "month" SMALLINT NOT NULL,
    "day" SMALLINT NOT NULL,
    "hour" SMALLINT NOT NULL,
    "minutes" SMALLINT NOT NULL,
    "seconds" SMALLINT NOT NULL,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "events";
-- +goose StatementEnd
