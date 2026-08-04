CREATE TABLE IF NOT EXISTS table_shard_keys (
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    table_name   TEXT NOT NULL,
    shard_key_column TEXT NOT NULL,
    is_manual    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (project_id, table_name)
);
