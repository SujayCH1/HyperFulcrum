import { NextResponse } from "next/server";
import pool from "@/lib/db";

export async function POST(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  try {
    const { rows } = await pool.query(
      `UPDATE shards
         SET schema_applied=true, last_executed_at=NOW()
       WHERE id=$1 RETURNING *`,
      [id]
    );
    if (!rows.length) return NextResponse.json({ error: "Shard not found" }, { status: 404 });
    const shard = rows[0];
    await pool.query(
      `INSERT INTO orchestration_logs (project_id, message, level)
       VALUES ($1, $2, 'cmd')`,
      [shard.project_id, `Schema applied to ${shard.shard_key} (${shard.node_host})`]
    );
    return NextResponse.json(shard);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Execution failed" }, { status: 500 });
  }
}
