CREATE TABLE shard_keys (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    table_name TEXT NOT NULL,
    key_column TEXT NOT NULL,

    updated_at TIMESTAMP NOT NULL,

    UNIQUE (project_id, table_name),
    CONSTRAINT fk_shard_keys_project
        FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);
