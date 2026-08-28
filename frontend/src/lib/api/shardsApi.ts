import apiClient from "@/lib/apiClient";
import type { Shard } from "@/types";

export interface CreateShardPayload {
  name: string;
  primary_node_id: string;
}

export const shardsApi = {
  list: async (projectId: string): Promise<Shard[]> => {
    const { data } = await apiClient.get<Shard[]>(`/projects/${projectId}/shards`);
    return Array.isArray(data) ? data : [];
  },
  get: async (shardId: string): Promise<Shard> => {
    const { data } = await apiClient.get<Shard>(`/shards/${shardId}`);
    return data;
  },
  create: async (projectId: string, payload: CreateShardPayload): Promise<Shard> => {
    const { data } = await apiClient.post<Shard>(`/projects/${projectId}/shards`, payload);
    return data;
  },
  rename: async (shardId: string, name: string): Promise<void> => {
    await apiClient.patch(`/shards/${shardId}/name`, null, { params: { name } });
  },
  remove: async (shardId: string): Promise<void> => {
    await apiClient.delete(`/shards/${shardId}`);
  },
};

