export interface Project {
  id: number;
  name: string;
  description?: string;
  status: "idle" | "active" | "error";
  schema_sql?: string;
  shard_count?: number;
  created_at: string;
}

export interface Shard {
  id: number;
  project_id: number;
  shard_key: string;
  node_host: string;
  status: "active" | "inactive" | "error";
  schema_applied: boolean;
  metrics_load: number;
  last_executed_at?: string;
  created_at: string;
}

export interface Log {
  id: number;
  project_id: number;
  message: string;
  level: "info" | "warn" | "error" | "cmd";
  timestamp: string;
}
