import { NextResponse } from "next/server";
import pool from "@/lib/db";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const projectId = searchParams.get("project_id");
  if (!projectId) return NextResponse.json({ error: "project_id required" }, { status: 400 });
  try {
    const { rows } = await pool.query(
      `SELECT * FROM orchestration_logs
       WHERE project_id=$1
       ORDER BY timestamp DESC LIMIT 200`,
      [projectId]
    );
    return NextResponse.json(rows);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to fetch logs" }, { status: 500 });
  }
}

export async function POST(request: Request) {
  try {
    const { project_id, message, level = "info" } = await request.json();
    const { rows } = await pool.query(
      `INSERT INTO orchestration_logs (project_id, message, level)
       VALUES ($1, $2, $3) RETURNING *`,
      [project_id, message, level]
    );
    return NextResponse.json(rows[0]);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to write log" }, { status: 500 });
  }
}
