import { NextResponse } from "next/server";
import pool from "@/lib/db";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const projectId = searchParams.get("project_id");
  if (!projectId) return NextResponse.json({ error: "project_id required" }, { status: 400 });
  try {
    const { rows } = await pool.query(
      "SELECT * FROM shards WHERE project_id=$1 ORDER BY created_at ASC",
      [projectId]
    );
    return NextResponse.json(rows);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to fetch shards" }, { status: 500 });
  }
}
