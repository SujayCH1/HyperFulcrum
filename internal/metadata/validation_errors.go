package metadata

import "errors"

var (
	ErrProjectRunning         = errors.New("project is running")
	ErrProjectHasNodes        = errors.New("project contains nodes")
	ErrProjectHasTopology     = errors.New("project contains replication topology")
	ErrProjectHasShards       = errors.New("project contains shards")
	ErrInvalidNodeRole        = errors.New("invalid node role")
	ErrDuplicateNodeName      = errors.New("node name already exists in project")
	ErrNodeActive             = errors.New("node is active")
	ErrNodeInTopology         = errors.New("node participates in replication topology")
	ErrNodeOwnsShard          = errors.New("node is the primary of a shard")
	ErrPrimaryNodeAlreadyUsed = errors.New("node is already primary of a shard")
	ErrDuplicateShardName     = errors.New("shard name already exists in project")
	ErrPrimaryRoleRequired    = errors.New("shard primary node must have the primary role")
	ErrShardHasStandbys       = errors.New("shard has standby relationships")
	ErrInvalidShardStatus     = errors.New("invalid shard status")
	ErrConnectionExists       = errors.New("node connection already exists")
	ErrConnectionNotFound     = errors.New("node connection not found")
	ErrInvalidConnection      = errors.New("invalid node connection")
	ErrTopologySelfRelation   = errors.New("primary and standby nodes must be different")
	ErrTopologyRoleMismatch   = errors.New("node roles do not match topology roles")
	ErrReplicaAlreadyUsed     = errors.New("standby node is already assigned")
	ErrShardIsReplica         = errors.New("primary node is already assigned as a standby")
	ErrDuplicateTopology      = errors.New("topology relation already exists")
	ErrSchemaNotLocked        = errors.New("project schema is not locked")
)
