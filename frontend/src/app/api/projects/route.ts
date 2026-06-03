import { NextResponse } from "next/server";
import { ensureTables } from "@/lib/dbInit";
import pool from "@/lib/db";

export async function GET() {
  await ensureTables();
  try {
    const { rows } = await pool.query(
      "SELECT * FROM projects ORDER BY created_at DESC"
    );
    return NextResponse.json(rows);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to fetch projects" }, { status: 500 });
  }
}

export async function POST(request: Request) {
  await ensureTables();
  try {
    const { name, description } = await request.json();
    const { rows } = await pool.query(
      `INSERT INTO projects (name, description, status)
       VALUES ($1, $2, 'idle') RETURNING *`,
      [name, description ?? ""]
    );
    const project = rows[0];
    await pool.query(
      `INSERT INTO orchestration_logs (project_id, message, level)
       VALUES ($1, 'Project created and ready for orchestration', 'info')`,
      [project.id]
    );
    return NextResponse.json(project);
  } catch (err) {
    console.error(err);
    return NextResponse.json({ error: "Failed to create project" }, { status: 500 });
  }
}
