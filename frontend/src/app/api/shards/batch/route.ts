import { NextResponse } from "next/server";
import pool from "@/lib/db";

export async function POST(request: Request) {
  try {
    const { project_id, count } = await request.json();
    if (!project_id || !count) {
      return NextResponse.json({ error: "project_id and count required" }, { status: 400 });
    }

    const client = await pool.connect();
    try {
      await client.query("BEGIN");
      const inserted: unknown[] = [];
      for (let i = 0; i < count; i++) {
        const shardKey = `shard_${String.fromCharCode(97 + i)}`;       // shard_a, shard_b, …
        const nodeHost = `node-${String(i + 1).padStart(2, "0")}.local`; // node-01.local, …
        const { rows } = await client.query(
          `INSERT INTO shards (project_id, shard_key, node_host, status)
           VALUES ($1, $2, $3, 'active') RETURNING *`,
          [project_id, shardKey, nodeHost]
        );
        inserted.push(rows[0]);
      }
      await client.query(
        `UPDATE projects SET shard_count=$1, status='active' WHERE id=$2`,
        [count, project_id]
      );
      await client.query(
        `INSERT INTO orchestration_logs (project_id, message, level)
         VALUES ($1, $2, 'info')`,
        [project_id, `Deployed ${count} shard node(s)`]
      );
      await client.query("COMMIT");
      return NextResponse.json(inserted);
    } catch (err) {
      await client.query("ROLLBACK");
      throw err;
    } finally {
      client.release();
    }
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to deploy shards" }, { status: 500 });
  }
}
