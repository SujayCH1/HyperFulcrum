// src/lib/api/projectsApi.ts
// All project-related calls to the Go backend.
// Endpoints (from internal/api/router/router.go):
//   GET    /projects/       → list all projects
//   POST   /projects        → create project  { name, description }
//   GET    /projects/ready  → list ready projects
//   GET    /projects/{id}   → get project by id
//   DELETE /projects/{id}   → delete project

import apiClient from "@/lib/apiClient";
import type { Project } from "@/types";

export const projectsApi = {
  /** GET /projects/ */
  list: async (): Promise<Project[]> => {
    const { data } = await apiClient.get<Project[]>("/projects/");
    return Array.isArray(data) ? data : [];
  },

  /** POST /projects  body: { name, description } */
  create: async (name: string, description: string): Promise<Project> => {
    const { data } = await apiClient.post<Project>("/projects", {
      name,
      description,
    });
    return data;
  },

  /** GET /projects/ready */
  listReady: async (): Promise<Project[]> => {
    const { data } = await apiClient.get<Project[]>("/projects/ready");
    return Array.isArray(data) ? data : [];
  },

  /** GET /projects/{id} */
  getById: async (id: string): Promise<Project> => {
    const { data } = await apiClient.get<Project>(`/projects/${id}`);
    return data;
  },

  /** DELETE /projects/{id} */
  remove: async (id: string): Promise<void> => {
    await apiClient.delete(`/projects/${id}`);
  },
};
