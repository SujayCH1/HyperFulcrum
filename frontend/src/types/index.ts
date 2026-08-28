// types/index.ts

export type NodeRole = "primary" | "standby" | "unassigned";
export type LogLevel = "info" | "error" | "warn" | "cmd";

// Maps to Go's project response shape
export interface Project {
  id: string;             // UUID string — was number, Go uses UUID
  name: string;
  description: string;
  node_count: number;
  ready: boolean;         // replaces old status string
  running: boolean;
  created_at: string;
  updated_at: string;
  schema_sql?: string;    // frontend-only, never sent to Go
}

// Maps to Go's node response shape
export interface Node {
  id: string;             // UUID string
  project_id: string;
  node_name: string;      // was: shard_key
  node_index: number;
  node_status: boolean;   // was: status / schema_applied
  role: NodeRole;
  created_at: string;
}

// Alias — keeps any existing `Shard` references compiling during migration
export interface Shard {
  id: string;
  project_id: string;
  shard_name: string;
  shard_index: number;
  primary_node_id: string;
  status: "provisioning" | "active" | "reconfiguring" | "unavailable";
  topology_generation: number;
  created_at: string;
  updated_at: string;
}

export interface NodeRuntimeState {
  node_id: string;
  observed_role: "primary" | "standby" | "unknown";
  postgres_status: "running" | "stopped" | "starting" | "bootstrapping" | "unreachable" | "unknown";
  postgres_version?: string;
  system_identifier?: string;
  timeline_id?: number;
  in_recovery?: boolean;
  read_only?: boolean;
  replication_lag_bytes?: number;
  last_observed_at?: string;
  observation_generation: number;
  last_error_code?: string;
  last_error_message?: string;
  updated_at: string;
}

// In-memory only — no backend storage
export interface Log {
  id: string;
  project_id: string;
  message: string;
  level: LogLevel;
  timestamp: string;
}

// Payload type for addNode / POST /projects/:id/nodes
export interface AddNodePayload {
  node_name: string;
  host: string;
  port: number;
  database_name: string;
  username: string;
  password: string;
}
