// // app/api/projects/route.ts
// // Go wraps all responses: { success: bool, message: string, data: any }
// // We unwrap .data before returning to the frontend store

// const GO_API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// export async function GET() {
//   try {
//     const res = await fetch(`${GO_API}/projects/`);
//     const json = await res.json();
//     if (!res.ok) return Response.json({ error: json.message }, { status: res.status });
//     // Unwrap Go's { success, message, data: Project[] }
//     return Response.json(json.data ?? []);
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Failed to fetch projects" }, { status: 500 });
//   }
// }

// export async function POST(request: Request) {
//   try {
//     const { name, description } = await request.json();
//     const res = await fetch(`${GO_API}/projects`, {
//       method: "POST",
//       headers: { "Content-Type": "application/json" },
//       body: JSON.stringify({ name, description }),
//     });
//     const json = await res.json();
//     if (!res.ok) return Response.json({ error: json.message }, { status: res.status });
//     // Unwrap Go's { success, message, data: Project }
//     return Response.json(json.data ?? json);
//   } catch (err) {
//     console.error(err);
//     return Response.json({ error: "Failed to create project" }, { status: 500 });
//   }
// }