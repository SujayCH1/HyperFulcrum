// // app/api/shards/[id]/execute/route.ts
// // Old: UPDATE shards SET schema_applied=true + INSERT log into DB
// // New: PATCH /nodes/:id/status?status=true on Go backend
// //      Log is written client-side via store.addLog()

// const GO_API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// // POST /api/shards/:id/execute → Go: PATCH /nodes/:id/status?status=true
// export async function POST(
//   _req: Request,
//   { params }: { params: Promise<{ id: string }> }
// ) {
//   const { id } = await params;
//   try {
//     const res = await fetch(`${GO_API}/nodes/${id}/status?status=true`, {
//       method: "PATCH",
//     });

//     if (!res.ok) {
//       const data = await res.json().catch(() => ({}));
//       return Response.json(
//         { error: data?.error ?? "Execution failed" },
//         { status: res.status }
//       );
//     }

//     // Go PATCH /nodes/:id/status returns 204 No Content on success
//     // Return a shape the caller recognises
//     return Response.json({ id, node_status: true, executed_at: new Date().toISOString() });
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Execution failed" }, { status: 500 });
//   }
// }
