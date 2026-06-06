// types/index.ts

export type NodeType = "shard" | "replica";
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
  node_type: NodeType;
  created_at: string;
}

// Alias — keeps any existing `Shard` references compiling during migration
export type Shard = Node;

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