CREATE TABLE IF NOT EXISTS node_connections (
    node_id UUID PRIMARY KEY,
    
    host TEXT NOT NULL,
    port INTEGER NOT NULL,
    database_name TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_node_connections_node
        FOREIGN KEY (node_id)
        REFERENCES nodes(id)
        ON DELETE CASCADE
);