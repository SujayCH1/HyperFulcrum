import pool from "./db";

let initialized = false;

export async function ensureTables(): Promise<void> {
  if (initialized) return;
  await pool.query(`
    CREATE TABLE IF NOT EXISTS projects (
      id        BIGSERIAL PRIMARY KEY,
      name      VARCHAR(255) NOT NULL,
      description TEXT,
      status    VARCHAR(50)  DEFAULT 'idle',
      schema_sql TEXT,
      shard_count INTEGER     DEFAULT 0,
      created_at TIMESTAMPTZ DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS shards (
      id            BIGSERIAL PRIMARY KEY,
      project_id    BIGINT REFERENCES projects(id) ON DELETE CASCADE,
      shard_key     VARCHAR(100) NOT NULL,
      node_host     VARCHAR(255) NOT NULL,
      status        VARCHAR(50)  DEFAULT 'active',
      schema_applied BOOLEAN     DEFAULT FALSE,
      metrics_load  INTEGER      DEFAULT 0,
      last_executed_at TIMESTAMPTZ,
      created_at    TIMESTAMPTZ  DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS orchestration_logs (
      id         BIGSERIAL PRIMARY KEY,
      project_id BIGINT REFERENCES projects(id) ON DELETE CASCADE,
      message    TEXT    NOT NULL,
      level      VARCHAR(20) DEFAULT 'info',
      timestamp  TIMESTAMPTZ DEFAULT NOW()
    );
  `);
  initialized = true;
}
