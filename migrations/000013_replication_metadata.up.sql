ALTER TYPE ntype RENAME VALUE 'shard' TO 'primary';
ALTER TYPE ntype RENAME VALUE 'replica' TO 'standby';
ALTER TYPE ntype ADD VALUE 'unassigned';

CREATE TYPE shard_status AS ENUM ('provisioning', 'active', 'reconfiguring', 'unavailable');

CREATE TABLE shards (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    shard_name VARCHAR(255) NOT NULL,
    shard_index INTEGER NOT NULL,
    primary_node_id UUID NOT NULL,
    status shard_status NOT NULL DEFAULT 'active',
    topology_generation BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT shards_project_id_id_unique UNIQUE (project_id, id),
    CONSTRAINT shards_project_id_id_primary_unique UNIQUE (project_id, id, primary_node_id),
    CONSTRAINT shards_project_name_unique UNIQUE (project_id, shard_name),
    CONSTRAINT shards_project_index_unique UNIQUE (project_id, shard_index),
    CONSTRAINT shards_primary_node_unique UNIQUE (primary_node_id),
    CONSTRAINT shards_index_nonnegative CHECK (shard_index >= 0),
    CONSTRAINT shards_topology_generation_positive CHECK (topology_generation > 0),
    CONSTRAINT fk_shards_project FOREIGN KEY (project_id)
        REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_shards_primary_node FOREIGN KEY (project_id, primary_node_id)
        REFERENCES nodes(project_id, id)
);

INSERT INTO shards (
    id, project_id, shard_name, shard_index, primary_node_id,
    status, created_at, updated_at
)
SELECT
    md5(id::TEXT || ':logical-shard')::UUID,
    project_id,
    node_name,
    node_index,
    id,
    CASE WHEN node_status THEN 'active'::shard_status
         ELSE 'unavailable'::shard_status END,
    created_at::TIMESTAMPTZ,
    NOW()
FROM nodes
WHERE node_type = 'primary';

CREATE TYPE topology_relationship_status AS ENUM
    ('planned', 'attaching', 'active', 'detaching', 'failed');

ALTER TABLE node_topology
ADD COLUMN shard_id UUID,
ADD COLUMN relationship_status topology_relationship_status NOT NULL DEFAULT 'planned',
ADD COLUMN replication_slot_name TEXT,
ADD COLUMN application_name TEXT,
ADD COLUMN promotion_priority INTEGER NOT NULL DEFAULT 0,
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE node_topology AS topology
SET shard_id = shard.id
FROM shards AS shard
WHERE shard.project_id = topology.project_id
  AND shard.primary_node_id = topology.shard_node_id;

ALTER TABLE node_topology
ALTER COLUMN shard_id SET NOT NULL,
ADD CONSTRAINT node_topology_shard_replica_unique UNIQUE (shard_id, replica_node_id),
ADD CONSTRAINT node_topology_promotion_priority_nonnegative CHECK (promotion_priority >= 0),
ADD CONSTRAINT fk_node_topology_logical_shard
    FOREIGN KEY (project_id, shard_id)
    REFERENCES shards(project_id, id),
ADD CONSTRAINT fk_node_topology_logical_shard_primary
    FOREIGN KEY (project_id, shard_id, shard_node_id)
    REFERENCES shards(project_id, id, primary_node_id)
    DEFERRABLE INITIALLY DEFERRED;

CREATE TYPE observed_node_role AS ENUM ('primary', 'standby', 'unknown');
CREATE TYPE postgres_runtime_status AS ENUM
    ('running', 'stopped', 'starting', 'bootstrapping', 'unreachable', 'unknown');

CREATE TABLE node_runtime_state (
    node_id UUID PRIMARY KEY,
    observed_role observed_node_role NOT NULL DEFAULT 'unknown',
    postgres_status postgres_runtime_status NOT NULL DEFAULT 'unknown',
    postgres_version TEXT,
    postgres_major_version INTEGER,
    system_identifier TEXT,
    timeline_id BIGINT,
    in_recovery BOOLEAN,
    read_only BOOLEAN,
    receive_lsn TEXT,
    replay_lsn TEXT,
    replication_lag_bytes BIGINT,
    last_agent_id UUID,
    last_observed_at TIMESTAMPTZ,
    observation_generation BIGINT NOT NULL DEFAULT 0,
    last_error_code TEXT,
    last_error_message TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT node_runtime_postgres_major_positive
        CHECK (postgres_major_version IS NULL OR postgres_major_version > 0),
    CONSTRAINT node_runtime_timeline_nonnegative
        CHECK (timeline_id IS NULL OR timeline_id >= 0),
    CONSTRAINT node_runtime_lag_nonnegative
        CHECK (replication_lag_bytes IS NULL OR replication_lag_bytes >= 0),
    CONSTRAINT node_runtime_observation_generation_nonnegative
        CHECK (observation_generation >= 0),
    CONSTRAINT fk_node_runtime_state_node FOREIGN KEY (node_id)
        REFERENCES nodes(id) ON DELETE CASCADE
);
