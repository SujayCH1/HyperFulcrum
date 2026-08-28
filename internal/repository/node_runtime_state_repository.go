package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	ObservedRolePrimary = "primary"
	ObservedRoleStandby = "standby"
	ObservedRoleUnknown = "unknown"

	PostgresStatusRunning       = "running"
	PostgresStatusStopped       = "stopped"
	PostgresStatusStarting      = "starting"
	PostgresStatusBootstrapping = "bootstrapping"
	PostgresStatusUnreachable   = "unreachable"
	PostgresStatusUnknown       = "unknown"
)

type NodeRuntimeState struct {
	NodeID                string     `json:"node_id"`
	ObservedRole          string     `json:"observed_role"`
	PostgresStatus        string     `json:"postgres_status"`
	PostgresVersion       *string    `json:"postgres_version,omitempty"`
	PostgresMajorVersion  *int       `json:"postgres_major_version,omitempty"`
	SystemIdentifier      *string    `json:"system_identifier,omitempty"`
	TimelineID            *int64     `json:"timeline_id,omitempty"`
	InRecovery            *bool      `json:"in_recovery,omitempty"`
	ReadOnly              *bool      `json:"read_only,omitempty"`
	ReceiveLSN            *string    `json:"receive_lsn,omitempty"`
	ReplayLSN             *string    `json:"replay_lsn,omitempty"`
	ReplicationLagBytes   *int64     `json:"replication_lag_bytes,omitempty"`
	LastAgentID           *string    `json:"last_agent_id,omitempty"`
	LastObservedAt        *time.Time `json:"last_observed_at,omitempty"`
	ObservationGeneration int64      `json:"observation_generation"`
	LastErrorCode         *string    `json:"last_error_code,omitempty"`
	LastErrorMessage      *string    `json:"last_error_message,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type NodeRuntimeStateRepository struct{ conn *sql.DB }

func NewNodeRuntimeStateRepository(conn *sql.DB) *NodeRuntimeStateRepository {
	return &NodeRuntimeStateRepository{conn: conn}
}

func (r *NodeRuntimeStateRepository) RuntimeStateGetByNodeID(ctx context.Context, nodeID string) (NodeRuntimeState, error) {
	return scanRuntimeState(r.conn.QueryRowContext(ctx, runtimeStateSelect+` WHERE state.node_id = $1`, nodeID))
}

func (r *NodeRuntimeStateRepository) RuntimeStateListByProject(ctx context.Context, projectID string) ([]NodeRuntimeState, error) {
	rows, err := r.conn.QueryContext(ctx, runtimeStateSelect+`
		JOIN nodes AS node ON node.id = state.node_id
		WHERE node.project_id = $1 ORDER BY node.node_index`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	states := make([]NodeRuntimeState, 0)
	for rows.Next() {
		state, err := scanRuntimeState(rows)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (r *NodeRuntimeStateRepository) RuntimeStateUpsert(ctx context.Context, state NodeRuntimeState) (bool, error) {
	if err := validateRuntimeState(state); err != nil {
		return false, err
	}
	result, err := r.conn.ExecContext(ctx, `
		INSERT INTO node_runtime_state (
			node_id, observed_role, postgres_status, postgres_version, postgres_major_version,
			system_identifier, timeline_id, in_recovery, read_only, receive_lsn, replay_lsn,
			replication_lag_bytes, last_agent_id, last_observed_at, observation_generation,
			last_error_code, last_error_message, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,NOW())
		ON CONFLICT (node_id) DO UPDATE SET
			observed_role = EXCLUDED.observed_role,
			postgres_status = EXCLUDED.postgres_status,
			postgres_version = EXCLUDED.postgres_version,
			postgres_major_version = EXCLUDED.postgres_major_version,
			system_identifier = EXCLUDED.system_identifier,
			timeline_id = EXCLUDED.timeline_id,
			in_recovery = EXCLUDED.in_recovery,
			read_only = EXCLUDED.read_only,
			receive_lsn = EXCLUDED.receive_lsn,
			replay_lsn = EXCLUDED.replay_lsn,
			replication_lag_bytes = EXCLUDED.replication_lag_bytes,
			last_agent_id = EXCLUDED.last_agent_id,
			last_observed_at = EXCLUDED.last_observed_at,
			observation_generation = EXCLUDED.observation_generation,
			last_error_code = EXCLUDED.last_error_code,
			last_error_message = EXCLUDED.last_error_message,
			updated_at = NOW()
		WHERE EXCLUDED.observation_generation > node_runtime_state.observation_generation
		   OR (EXCLUDED.observation_generation = node_runtime_state.observation_generation
		       AND EXCLUDED.last_observed_at > COALESCE(
		           node_runtime_state.last_observed_at, '-infinity'::TIMESTAMPTZ))`,
		state.NodeID, state.ObservedRole, state.PostgresStatus, state.PostgresVersion,
		state.PostgresMajorVersion, state.SystemIdentifier, state.TimelineID, state.InRecovery,
		state.ReadOnly, state.ReceiveLSN, state.ReplayLSN, state.ReplicationLagBytes,
		state.LastAgentID, state.LastObservedAt, state.ObservationGeneration,
		state.LastErrorCode, state.LastErrorMessage)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}

func (r *NodeRuntimeStateRepository) RuntimeStateMarkUnreachable(
	ctx context.Context, nodeID string, observedAt time.Time, generation int64,
	errorCode, errorMessage string,
) (bool, error) {
	result, err := r.conn.ExecContext(ctx, `
		INSERT INTO node_runtime_state (
			node_id, observed_role, postgres_status, last_observed_at,
			observation_generation, last_error_code, last_error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (node_id) DO UPDATE SET
			postgres_status = EXCLUDED.postgres_status,
			last_observed_at = EXCLUDED.last_observed_at,
			observation_generation = EXCLUDED.observation_generation,
			last_error_code = EXCLUDED.last_error_code,
			last_error_message = EXCLUDED.last_error_message,
			updated_at = NOW()
		WHERE EXCLUDED.observation_generation > node_runtime_state.observation_generation
		   OR (EXCLUDED.observation_generation = node_runtime_state.observation_generation
		       AND EXCLUDED.last_observed_at > COALESCE(
		           node_runtime_state.last_observed_at, '-infinity'::TIMESTAMPTZ))`,
		nodeID, ObservedRoleUnknown, PostgresStatusUnreachable, observedAt, generation,
		errorCode, errorMessage)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}

func (r *NodeRuntimeStateRepository) RuntimeStateRemove(ctx context.Context, nodeID string) error {
	return requireAffected(r.conn.ExecContext(ctx, `DELETE FROM node_runtime_state WHERE node_id = $1`, nodeID))
}

const runtimeStateSelect = `SELECT
	state.node_id, state.observed_role, state.postgres_status, state.postgres_version,
	state.postgres_major_version, state.system_identifier, state.timeline_id,
	state.in_recovery, state.read_only, state.receive_lsn, state.replay_lsn,
	state.replication_lag_bytes, state.last_agent_id, state.last_observed_at,
	state.observation_generation, state.last_error_code, state.last_error_message,
	state.updated_at FROM node_runtime_state AS state`

func scanRuntimeState(row rowScanner) (NodeRuntimeState, error) {
	var state NodeRuntimeState
	err := row.Scan(&state.NodeID, &state.ObservedRole, &state.PostgresStatus,
		&state.PostgresVersion, &state.PostgresMajorVersion, &state.SystemIdentifier,
		&state.TimelineID, &state.InRecovery, &state.ReadOnly, &state.ReceiveLSN,
		&state.ReplayLSN, &state.ReplicationLagBytes, &state.LastAgentID,
		&state.LastObservedAt, &state.ObservationGeneration, &state.LastErrorCode,
		&state.LastErrorMessage, &state.UpdatedAt)
	return state, err
}

func validateRuntimeState(state NodeRuntimeState) error {
	switch state.ObservedRole {
	case ObservedRolePrimary, ObservedRoleStandby, ObservedRoleUnknown:
	default:
		return fmt.Errorf("invalid observed node role: %s", state.ObservedRole)
	}
	switch state.PostgresStatus {
	case PostgresStatusRunning, PostgresStatusStopped, PostgresStatusStarting,
		PostgresStatusBootstrapping, PostgresStatusUnreachable, PostgresStatusUnknown:
		return nil
	default:
		return fmt.Errorf("invalid PostgreSQL runtime status: %s", state.PostgresStatus)
	}
}
