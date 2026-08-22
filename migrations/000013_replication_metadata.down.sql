-- Rolling back is intentionally refused while an unassigned node exists,
-- because the legacy enum has no truthful value for that role.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM nodes WHERE node_type = 'unassigned') THEN
        RAISE EXCEPTION 'cannot roll back replication metadata while unassigned nodes exist';
    END IF;
END
$$;

DROP TABLE IF EXISTS node_runtime_state;
DROP TYPE IF EXISTS postgres_runtime_status;
DROP TYPE IF EXISTS observed_node_role;

ALTER TABLE node_topology
DROP CONSTRAINT IF EXISTS fk_node_topology_shard_primary,
DROP CONSTRAINT IF EXISTS fk_node_topology_shard,
DROP CONSTRAINT IF EXISTS node_topology_promotion_priority_nonnegative,
DROP CONSTRAINT IF EXISTS node_topology_shard_replica_unique,
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS promotion_priority,
DROP COLUMN IF EXISTS application_name,
DROP COLUMN IF EXISTS replication_slot_name,
DROP COLUMN IF EXISTS relationship_status,
DROP COLUMN IF EXISTS shard_id;

DROP TYPE IF EXISTS topology_relationship_status;
DROP TABLE IF EXISTS shards;
DROP TYPE IF EXISTS shard_status;

-- PostgreSQL enum values cannot be removed directly, so recreate the legacy
-- type after proving that no value would be lost.
ALTER TABLE nodes ALTER COLUMN node_type DROP DEFAULT;
ALTER TYPE ntype RENAME TO ntype_replication;
CREATE TYPE ntype AS ENUM ('shard', 'replica');
ALTER TABLE nodes
ALTER COLUMN node_type TYPE ntype
USING (
    CASE node_type::TEXT
        WHEN 'primary' THEN 'shard'
        WHEN 'standby' THEN 'replica'
    END
)::ntype;
DROP TYPE ntype_replication;
