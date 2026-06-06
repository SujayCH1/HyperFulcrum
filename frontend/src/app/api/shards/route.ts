// // app/api/shards/route.ts
// // Go wraps: { success, message, data: Node[] }

// const GO_API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// export async function GET(request: Request) {
//   const { searchParams } = new URL(request.url);
//   const projectId = searchParams.get("project_id");

//   if (!projectId) {
//     return Response.json({ error: "project_id required" }, { status: 400 });
//   }

//   try {
//     const res = await fetch(`${GO_API}/projects/${projectId}/nodes`);
//     const json = await res.json();
//     if (!res.ok) return Response.json({ error: json.message }, { status: res.status });
//     // Unwrap Go's { success, message, data: Node[] }
//     return Response.json(json.data ?? []);
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Failed to fetch nodes" }, { status: 500 });
//   }
// }