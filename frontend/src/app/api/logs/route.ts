// // app/api/logs/route.ts
// // No DB, no Go backend endpoint for logs yet.
// // Logs are managed in-memory via Zustand (projectStore → addLog / clearLogs).
// // These endpoints are kept as stubs so any existing fetch("/api/logs") calls
// // don't 404 — they just return the in-memory log array passed in the request.

// import { NextResponse } from "next/server";

// // GET /api/logs?project_id=X
// // Returns empty array — caller should read logs from Zustand store directly
// // instead of fetching this route. This stub prevents 404s during migration.
// export async function GET(request: Request) {
//   const { searchParams } = new URL(request.url);
//   const projectId = searchParams.get("project_id");

//   if (!projectId) {
//     return NextResponse.json({ error: "project_id required" }, { status: 400 });
//   }

//   // No DB — logs live in Zustand store on the client
//   return NextResponse.json([]);
// }

// // POST /api/logs
// // No-op stub — logging is done via store.addLog() on the client side.
// export async function POST(request: Request) {
//   try {
//     const body = await request.json();
//     // Acknowledge the write without persisting anything
//     return NextResponse.json({
//       id: crypto.randomUUID(),
//       project_id: body.project_id ?? "",
//       message: body.message ?? "",
//       level: body.level ?? "info",
//       timestamp: new Date().toISOString(),
//     });
//   } catch {
//     return NextResponse.json({ error: "Invalid request body" }, { status: 400 });
//   }
// }