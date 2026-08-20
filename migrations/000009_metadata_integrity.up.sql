ALTER TABLE projects
ADD CONSTRAINT projects_node_count_nonnegative
CHECK (node_count >= 0);

ALTER TABLE nodes
ADD CONSTRAINT nodes_project_id_id_unique
UNIQUE (project_id, id),
ADD CONSTRAINT nodes_project_id_name_unique
UNIQUE (project_id, node_name),
ADD CONSTRAINT nodes_project_id_index_unique
UNIQUE (project_id, node_index),
ADD CONSTRAINT nodes_index_nonnegative
CHECK (node_index >= 0),
ADD CONSTRAINT fk_nodes_project
FOREIGN KEY (project_id)
REFERENCES projects(id)
ON DELETE CASCADE;

ALTER TABLE node_connections
ADD CONSTRAINT node_connections_valid_port
CHECK (port BETWEEN 1 AND 65535);

ALTER TABLE node_topology
ADD CONSTRAINT node_topology_distinct_nodes
CHECK (shard_node_id <> replica_node_id),
ADD CONSTRAINT node_topology_replica_unique
UNIQUE (replica_node_id),
ADD CONSTRAINT fk_node_topology_project
FOREIGN KEY (project_id)
REFERENCES projects(id)
ON DELETE CASCADE,
ADD CONSTRAINT fk_node_topology_shard
FOREIGN KEY (project_id, shard_node_id)
REFERENCES nodes(project_id, id)
ON DELETE CASCADE,
ADD CONSTRAINT fk_node_topology_replica
FOREIGN KEY (project_id, replica_node_id)
REFERENCES nodes(project_id, id)
ON DELETE CASCADE;

ALTER TABLE columns
ADD CONSTRAINT fk_columns_project
FOREIGN KEY (project_id)
REFERENCES projects(id)
ON DELETE CASCADE;

ALTER TABLE fk_edges
ADD CONSTRAINT fk_fk_edges_project
FOREIGN KEY (project_id)
REFERENCES projects(id)
ON DELETE CASCADE;

ALTER TABLE schema_versions
ADD CONSTRAINT schema_versions_project_unique
UNIQUE (project_id),
ADD CONSTRAINT fk_schema_versions_project
FOREIGN KEY (project_id)
REFERENCES projects(id)
ON DELETE CASCADE;
