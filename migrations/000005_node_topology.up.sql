CREATE TABLE node_topology (
    relation_id UUID PRIMARY KEY,

    project_id UUID NOT NULL,
    shard_node_id UUID NOT NULL,
    replica_node_id UUID NOT NULL,

    created_at TIMESTAMP NOT NULL,

    CONSTRAINT unique_shard_replica
    UNIQUE (shard_node_id, replica_node_id)
);