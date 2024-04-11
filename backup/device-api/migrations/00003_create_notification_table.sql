-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_events (
                                              id          serial      PRIMARY KEY,
                                              device_id   bigint,
                                              message       TEXT,
                                              lang int2 NOT NULL,
                                              status int2 NOT NULL,
                                              payload     jsonb,
                                              created_at  timestamp   DEFAULT now() NOT NULL,
                                              updated_at  timestamp   DEFAULT now() NOT NULL,

    FOREIGN KEY (device_id) REFERENCES devices (id)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_events;
-- +goose StatementEnd
