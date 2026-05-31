CREATE TABLE nodes (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    node_name VARCHAR NOT NULL,
    node_index INTEGER NOT NULL,
    node_status BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
)