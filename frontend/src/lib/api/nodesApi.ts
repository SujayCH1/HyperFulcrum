// src/lib/api/nodesApi.ts
// All node-related calls to the Go backend.
// Endpoints (from internal/api/router/router.go):
//   POST   /projects/{projectId}/nodes          → add node
//   GET    /projects/{projectId}/nodes          → list nodes for project
//   DELETE /nodes/{id}                          → delete node
//   PUT    /nodes/{id}/name?name=<value>        → rename node
//   PATCH  /nodes/{id}/status?status=<bool>     → update node status
//   PATCH  /nodes/{id}/role?role=<value>        → update desired node role
//
// Go NodeDto fields: { name, role, index, status }

import apiClient from "@/lib/apiClient";
import type { Node, NodeRole, NodeRuntimeState } from "@/types";

export interface CreateNodePayload {
  name: string;       // Go field: "name"
  role: NodeRole;
  // index: number;      // Go field: "index"
  // status: boolean;    // Go field: "status"
}

export const nodesApi = {
  /** GET /projects/{projectId}/nodes */
  list: async (projectId: string): Promise<Node[]> => {
    const { data } = await apiClient.get<Node[]>(
      `/projects/${projectId}/nodes`
    );
    return Array.isArray(data) ? data : [];
  },

  /** POST /projects/{projectId}/nodes */
  create: async (
    projectId: string,
    payload: CreateNodePayload
  ): Promise<Node | null> => {
    const { data } = await apiClient.post<Node>(
      `/projects/${projectId}/nodes`,
      payload
    );
    return data;
  },

  /**
   * Batch create N nodes for a project.
   * Names are auto-generated: shard_a, shard_b, …
   */
  batchCreate: async (
    projectId: string,
    count: number
  ): Promise<{ inserted: (Node | null)[]; failed: { index: number; error: string }[] }> => {
    const inserted: (Node | null)[] = [];
    const failed: { index: number; error: string }[] = [];

    for (let i = 0; i < count; i++) {
      const nodeName = `shard_${String.fromCharCode(97 + i)}`; // shard_a, shard_b, …
      try {
        const node = await nodesApi.create(projectId, {
          name: nodeName,
          role: "primary",
          // index: i,
          // status: false,
        });
        inserted.push(node);
      } catch (err) {
        failed.push({ index: i, error: (err as Error).message });
      }
    }

    return { inserted, failed };
  },

  /** DELETE /nodes/{id} */
  remove: async (nodeId: string): Promise<void> => {
    await apiClient.delete(`/nodes/${nodeId}`);
  },

  /** PUT /nodes/{id}/name?name=<value> */
  updateName: async (nodeId: string, name: string): Promise<void> => {
    await apiClient.put(`/nodes/${nodeId}/name`, null, {
      params: { name },
    });
  },

  /** PATCH /nodes/{id}/status?status=<bool> */
  updateStatus: async (nodeId: string, status: boolean): Promise<void> => {
    await apiClient.patch(`/nodes/${nodeId}/status`, null, {
      params: { status },
    });
  },

  /** PATCH /nodes/{id}/role?role=<value> */
  updateRole: async (nodeId: string, role: NodeRole): Promise<void> => {
    await apiClient.patch(`/nodes/${nodeId}/role`, null, {
      params: { role },
    });
  },

  runtimeState: async (nodeId: string): Promise<NodeRuntimeState> => {
    const { data } = await apiClient.get<NodeRuntimeState>(`/nodes/${nodeId}/runtime-state`);
    return data;
  },
};
