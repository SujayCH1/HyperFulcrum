// // app/api/shards/batch/route.ts
// // Calls Go POST /projects/:id/nodes N times
// // Go NodeDto expects: { name, role, index, status }

// const GO_API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// export async function POST(request: Request) {
//   try {
//     const { project_id, count } = await request.json();

//     if (!project_id || !count) {
//       return Response.json({ error: "project_id and count required" }, { status: 400 });
//     }

//     const inserted: unknown[] = [];
//     const failed: { index: number; error: string }[] = [];

//     for (let i = 0; i < count; i++) {
//       const nodeName = `shard_${String.fromCharCode(97 + i)}`; // shard_a, shard_b, …

//       try {
//         const res = await fetch(`${GO_API}/projects/${project_id}/nodes`, {
//           method: "POST",
//           headers: { "Content-Type": "application/json" },
//           body: JSON.stringify({
//             name: nodeName,   // Go NodeDto field: "name"
//             role: "primary",
//             index: i,         // Go NodeDto field: "index"
//             status: false,    // Go NodeDto field: "status"
//           }),
//         });

//         const data = await res.json();

//         if (!res.ok) {
//           failed.push({ index: i, error: data?.message ?? data?.error ?? "Unknown error" });
//         } else {
//           inserted.push(data);
//         }
//       } catch (err) {
//         failed.push({ index: i, error: (err as Error).message });
//       }
//     }

//     if (failed.length > 0) {
//       return Response.json(
//         { inserted, failed, message: `${inserted.length} of ${count} nodes created` },
//         { status: failed.length === count ? 500 : 207 }
//       );
//     }

//     return Response.json(inserted);
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Failed to deploy nodes" }, { status: 500 });
//   }
// }
