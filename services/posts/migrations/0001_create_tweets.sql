CREATE TABLE IF NOT EXISTS tweets (
    id  BIGSERIAL   PRIMARY KEY,
    user_id BIGINT  NOT NULL,
    text    text    NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);