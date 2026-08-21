package metadata

import "errors"

var (
	ErrProjectRunning       = errors.New("project is running")
	ErrProjectHasNodes      = errors.New("project contains nodes")
	ErrProjectHasTopology   = errors.New("project contains replication topology")
	ErrInvalidNodeType      = errors.New("invalid node type")
	ErrDuplicateNodeName    = errors.New("node name already exists in project")
	ErrNodeActive           = errors.New("node is active")
	ErrNodeInTopology       = errors.New("node participates in replication topology")
	ErrConnectionExists     = errors.New("node connection already exists")
	ErrConnectionNotFound   = errors.New("node connection not found")
	ErrInvalidConnection    = errors.New("invalid node connection")
	ErrTopologySelfRelation = errors.New("shard and replica nodes must be different")
	ErrTopologyRoleMismatch = errors.New("node types do not match topology roles")
	ErrReplicaAlreadyUsed   = errors.New("replica node is already assigned")
	ErrShardIsReplica       = errors.New("shard node is already assigned as a replica")
	ErrDuplicateTopology    = errors.New("topology relation already exists")
)
