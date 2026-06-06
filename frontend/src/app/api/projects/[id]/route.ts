// // app/api/projects/[id]/route.ts

// const GO_API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// // GET /api/projects/:id → Go: GET /projects/:id
// export async function GET(
//   _req: Request,
//   { params }: { params: Promise<{ id: string }> }
// ) {
//   const { id } = await params;
//   try {
//     const res = await fetch(`${GO_API}/projects/${id}`);
//     const data = await res.json();
//     if (!res.ok) return Response.json(data, { status: res.status });
//     return Response.json(data);
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Failed to fetch project" }, { status: 500 });
//   }
// }

// // PUT /api/projects/:id
// // schema_sql, status, shard_count are frontend-only fields — Go doesn't store them.
// // This route acknowledges the update without hitting Go.
// // The real update happens in Zustand via store.updateProjectSchema().
// export async function PUT(
//   request: Request,
//   { params }: { params: Promise<{ id: string }> }
// ) {
//   const { id } = await params;
//   try {
//     const body = await request.json();

//     // Nothing to proxy — all updatable fields are frontend-only
//     // Return the same shape the caller expects
//     return Response.json({
//       id,
//       schema_sql: body.schema_sql ?? null,
//       status: body.status ?? null,
//       shard_count: body.shard_count ?? null,
//       updated_at: new Date().toISOString(),
//     });
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Invalid request body" }, { status: 400 });
//   }
// }

// // DELETE /api/projects/:id → Go: DELETE /projects/:id
// export async function DELETE(
//   _req: Request,
//   { params }: { params: Promise<{ id: string }> }
// ) {
//   const { id } = await params;
//   try {
//     const res = await fetch(`${GO_API}/projects/${id}`, { method: "DELETE" });
//     if (!res.ok) {
//       const data = await res.json().catch(() => ({}));
//       return Response.json(data, { status: res.status });
//     }
//     return Response.json({ ok: true });
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Failed to delete project" }, { status: 500 });
//   }
// }