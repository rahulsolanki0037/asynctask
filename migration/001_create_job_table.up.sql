CREATE TABLE IF NOT EXISTS jobs (
    id int PRIMARY KEY,
    type varchar(100) NOT NULL,
    payload TEXT NOT NULL,
    status varchar(25) NOT NULL,
    retry_count int8 NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);