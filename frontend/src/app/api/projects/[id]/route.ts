import { NextResponse } from "next/server";
import pool from "@/lib/db";

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  try {
    const { rows } = await pool.query("SELECT * FROM projects WHERE id=$1", [id]);
    if (!rows.length) return NextResponse.json({ error: "Not found" }, { status: 404 });
    return NextResponse.json(rows[0]);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed" }, { status: 500 });
  }
}

export async function PUT(
  request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  try {
    const body = await request.json();
    const fields: string[] = [];
    const vals: unknown[] = [];
    let idx = 1;
    if (body.schema_sql !== undefined) { fields.push(`schema_sql=$${idx++}`); vals.push(body.schema_sql); }
    if (body.status !== undefined)     { fields.push(`status=$${idx++}`);     vals.push(body.status); }
    if (body.shard_count !== undefined){ fields.push(`shard_count=$${idx++}`);vals.push(body.shard_count); }
    if (!fields.length) return NextResponse.json({ error: "Nothing to update" }, { status: 400 });
    vals.push(id);
    const { rows } = await pool.query(
      `UPDATE projects SET ${fields.join(", ")} WHERE id=$${idx} RETURNING *`,
      vals
    );
    return NextResponse.json(rows[0]);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to update" }, { status: 500 });
  }
}

export async function DELETE(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  try {
    await pool.query("DELETE FROM projects WHERE id=$1", [id]);
    return NextResponse.json({ ok: true });
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to delete" }, { status: 500 });
  }
}
