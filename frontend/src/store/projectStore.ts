// src/store/projectStore.ts
// Zustand store — all backend calls go directly to the Go API
// via the typed modules in @/lib/api (no Next.js proxy routes).

import { create } from "zustand";
import { projectsApi, nodesApi } from "@/lib/api";
import type { Project, Node, Log, AddNodePayload } from "@/types";

interface ProjectState {
  projects: Project[];
  activeProject: Project | null;
  nodes: Node[];
  logs: Log[];

  loadingProjects: boolean;
  loadingShards: boolean;
  loadingLogs: boolean;
  error: string | null;

  drawerOpen: boolean;
  activeTab: "schema" | "shards" | "logs";

  fetchProjects: () => Promise<void>;
  createProject: (name: string, description: string) => Promise<Project | null>;
  setActiveProject: (project: Project | null) => void;
  updateProjectSchema: (projectId: string, schemaSql: string) => void;

  fetchNodes: (projectId: string) => Promise<void>;
  addNode: (projectId: string, payload: AddNodePayload) => Promise<boolean>;
  deleteNode: (nodeId: string) => Promise<boolean>;
  updateNodeStatus: (nodeId: string, status: boolean) => Promise<boolean>;
  updateNodeName: (nodeId: string, name: string) => Promise<boolean>;
  updateNodeType: (nodeId: string, type: string) => Promise<boolean>;

  addLog: (message: string, level?: "info" | "error" | "warn" | "cmd") => void;
  clearLogs: () => void;

  openDrawer: () => void;
  closeDrawer: () => void;
  setActiveTab: (tab: "schema" | "shards" | "logs") => void;
  clearError: () => void;
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projects: [],
  activeProject: null,
  nodes: [],
  logs: [],
  loadingProjects: false,
  loadingShards: false,
  loadingLogs: false,
  error: null,
  drawerOpen: false,
  activeTab: "schema",

  // ── Projects ──────────────────────────────────────────────────────────────

  fetchProjects: async () => {
    set({ loadingProjects: true, error: null });
    try {
      const projects = await projectsApi.list();
      set({ projects });
    } catch (e) {
      set({ error: (e as Error).message });
    } finally {
      set({ loadingProjects: false });
    }
  },

  createProject: async (name, description) => {
    set({ error: null });
    try {
      const project = await projectsApi.create(name, description);
      set((s) => ({ projects: [project, ...s.projects] }));
      get().addLog(`Project "${project.name}" created`, "info");
      return project;
    } catch (e) {
      set({ error: (e as Error).message });
      return null;
    }
  },

  setActiveProject: (project) => {
    set({ activeProject: project, nodes: [], logs: [], activeTab: "schema" });
  },

  // schema_sql is frontend-only — never sent to Go
  updateProjectSchema: (projectId, schemaSql) => {
    set((s) => ({
      projects: s.projects.map((p) =>
        p.id === projectId ? { ...p, schema_sql: schemaSql } : p
      ),
      activeProject:
        s.activeProject?.id === projectId
          ? { ...s.activeProject, schema_sql: schemaSql }
          : s.activeProject,
    }));
  },

  // ── Nodes ─────────────────────────────────────────────────────────────────

  fetchNodes: async (projectId) => {
    set({ loadingShards: true, error: null });
    try {
      const nodes = await nodesApi.list(projectId);
      set({ nodes });
    } catch (e) {
      set({ error: (e as Error).message });
    } finally {
      set({ loadingShards: false });
    }
  },

  addNode: async (projectId, payload) => {
    try {
      await nodesApi.create(projectId, {
        name: payload.node_name,
        type: "shard",
        // index: 0,
        // status: false,
      });
      await get().fetchNodes(projectId);
      get().addLog(`Node "${payload.node_name}" added`, "info");
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      get().addLog(`Failed to add node: ${(e as Error).message}`, "error");
      return false;
    }
  },

  deleteNode: async (nodeId) => {
    try {
      await nodesApi.remove(nodeId);
      set((s) => ({ nodes: s.nodes.filter((n) => n.id !== nodeId) }));
      get().addLog(`Node ${nodeId} removed`, "warn");
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      return false;
    }
  },

  updateNodeStatus: async (nodeId, status) => {
    try {
      await nodesApi.updateStatus(nodeId, status);
      set((s) => ({
        nodes: s.nodes.map((n) =>
          n.id === nodeId ? { ...n, node_status: status } : n
        ),
      }));
      get().addLog(`Node ${nodeId} status → ${status}`, "cmd");
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      get().addLog(`Status update failed: ${(e as Error).message}`, "error");
      return false;
    }
  },

  updateNodeName: async (nodeId, name) => {
    try {
      await nodesApi.updateName(nodeId, name);
      set((s) => ({
        nodes: s.nodes.map((n) =>
          n.id === nodeId ? { ...n, node_name: name } : n
        ),
      }));
      get().addLog(`Node ${nodeId} renamed to "${name}"`, "info");
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      return false;
    }
  },

  updateNodeType: async (nodeId, type) => {
    try {
      await nodesApi.updateType(nodeId, type);
      set((s) => ({
        nodes: s.nodes.map((n) =>
          n.id === nodeId ? { ...n, node_type: type as Node["node_type"] } : n
        ),
      }));
      get().addLog(`Node ${nodeId} type → ${type}`, "info");
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      return false;
    }
  },

  // ── Logs (in-memory) ──────────────────────────────────────────────────────

  addLog: (message, level = "info") => {
    const log: Log = {
      id: crypto.randomUUID(),
      project_id: get().activeProject?.id ?? "",
      message,
      level,
      timestamp: new Date().toISOString(),
    };
    set((s) => ({ logs: [log, ...s.logs] }));
  },

  clearLogs: () => set({ logs: [] }),

  openDrawer: () => set({ drawerOpen: true }),
  closeDrawer: () => set({ drawerOpen: false }),
  setActiveTab: (tab) => set({ activeTab: tab }),
  clearError: () => set({ error: null }),
}));
