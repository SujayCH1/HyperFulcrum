// src/lib/api/index.ts
// Barrel export — import both API modules from one place:
//   import { projectsApi, nodesApi } from "@/lib/api";

export { projectsApi } from "./projectsApi";
export { nodesApi } from "./nodesApi";
export type { CreateNodePayload } from "./nodesApi";
