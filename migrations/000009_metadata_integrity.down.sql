ALTER TABLE schema_versions
DROP CONSTRAINT IF EXISTS fk_schema_versions_project,
DROP CONSTRAINT IF EXISTS schema_versions_project_unique;

ALTER TABLE fk_edges
DROP CONSTRAINT IF EXISTS fk_fk_edges_project;

ALTER TABLE columns
DROP CONSTRAINT IF EXISTS fk_columns_project;

ALTER TABLE node_topology
DROP CONSTRAINT IF EXISTS fk_node_topology_replica,
DROP CONSTRAINT IF EXISTS fk_node_topology_shard,
DROP CONSTRAINT IF EXISTS fk_node_topology_project,
DROP CONSTRAINT IF EXISTS node_topology_replica_unique,
DROP CONSTRAINT IF EXISTS node_topology_distinct_nodes;

ALTER TABLE node_connections
DROP CONSTRAINT IF EXISTS node_connections_valid_port;

ALTER TABLE nodes
DROP CONSTRAINT IF EXISTS fk_nodes_project,
DROP CONSTRAINT IF EXISTS nodes_index_nonnegative,
DROP CONSTRAINT IF EXISTS nodes_project_id_index_unique,
DROP CONSTRAINT IF EXISTS nodes_project_id_name_unique,
DROP CONSTRAINT IF EXISTS nodes_project_id_id_unique;

ALTER TABLE projects
DROP CONSTRAINT IF EXISTS projects_node_count_nonnegative;
