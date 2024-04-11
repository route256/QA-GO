-- +goose Up
CREATE TABLE IF NOT EXISTS devices_events (
    id          serial      PRIMARY KEY,
    device_id   bigint,
    type        int2        NOT NULL,
    status      int2        NOT NULL,
    payload     jsonb,
    created_at  timestamp   DEFAULT now() NOT NULL,
    updated_at  timestamp   DEFAULT now() NOT NULL,

    FOREIGN KEY (device_id) REFERENCES devices (id)
);

INSERT INTO devices_events (device_id, type, status, payload, created_at, updated_at)
VALUES  (1, 1, 1, '{"user_id": 24154, "platform": "Ios", "entered_at": "2021-11-09T13:51:49.870041258Z"}', '2021-11-09 13:51:49.876400', '2021-11-09 13:51:49.876400'),
        (2, 1, 1, '{"user_id": 28123, "platform": "Android", "entered_at": "2021-11-09T13:52:16.785666915Z"}', '2021-11-09 13:52:16.791872', '2021-11-09 13:52:16.791872'),
        (3, 1, 1, '{"user_id": 412414, "platform": "Android", "entered_at": "2021-11-09T13:52:20.094174185Z"}', '2021-11-09 13:52:20.100511', '2021-11-09 13:52:20.100511'),
        (4, 1, 1, '{"user_id": 41244, "platform": "Linux", "entered_at": "2021-11-09T13:52:26.945600822Z"}', '2021-11-09 13:52:26.951710', '2021-11-09 13:52:26.951710'),
        (5, 1, 1, '{"user_id": 76262, "platform": "Windows", "entered_at": "2021-11-09T13:52:45.590365311Z"}', '2021-11-09 13:52:45.596453', '2021-11-09 13:52:45.596453'),
        (3, 1, 1, 'null', '2021-11-09 13:52:56.885199', '2021-11-09 13:52:56.885199'),
        (6, 1, 1, '{"user_id": 52352, "platform": "Windows", "entered_at": "2021-11-09T13:53:07.152502881Z"}', '2021-11-09 13:53:07.158634', '2021-11-09 13:53:07.158634'),
        (5, 1, 1, 'null', '2021-11-09 13:53:13.451880', '2021-11-09 13:53:13.451880'),
        (7, 1, 1, '{"user_id": 241442, "platform": "Ios", "entered_at": "2021-11-09T13:53:26.35022303Z"}', '2021-11-09 13:53:26.356291', '2021-11-09 13:53:26.356291'),
        (8, 1, 1, '{"user_id": 15515, "platform": "Ios", "entered_at": "2021-11-09T13:53:29.226157645Z"}', '2021-11-09 13:53:29.231764', '2021-11-09 13:53:29.231764'),
        (9, 1, 1, '{"user_id": 74349, "platform": "Android", "entered_at": "2021-11-09T13:53:39.733289098Z"}', '2021-11-09 13:53:39.739527', '2021-11-09 13:53:39.739527'),
        (10, 1, 1, '{"user_id": 124141, "platform": "Linux", "entered_at": "2021-11-09T13:54:04.246343598Z"}', '2021-11-09 13:54:04.252218', '2021-11-09 13:54:04.252218');

-- +goose Down
DROP TABLE IF EXISTS devices_events;
