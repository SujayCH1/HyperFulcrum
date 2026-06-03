import { create } from "zustand";
import type { Project, Shard, Log } from "@/types";

// ─── Types ───────────────────────────────────────────────────────────────────

interface ProjectState {
  // ── Data
  projects: Project[];
  activeProject: Project | null;
  shards: Shard[];
  logs: Log[];

  // ── Loading / error flags
  loadingProjects: boolean;
  loadingShards: boolean;
  loadingLogs: boolean;
  error: string | null;

  // ── UI flags
  drawerOpen: boolean;
  activeTab: "schema" | "shards" | "logs";

  // ── Actions: projects
  fetchProjects: () => Promise<void>;
  createProject: (name: string, description: string) => Promise<Project | null>;
  setActiveProject: (project: Project | null) => void;
  updateProjectSchema: (projectId: number, schemaSql: string) => Promise<boolean>;

  // ── Actions: shards
  fetchShards: (projectId: number) => Promise<void>;
  deployShards: (projectId: number, count: number) => Promise<boolean>;
  executeShard: (shardId: number) => Promise<boolean>;

  // ── Actions: logs
  fetchLogs: (projectId: number) => Promise<void>;

  // ── Actions: UI
  openDrawer: () => void;
  closeDrawer: () => void;
  setActiveTab: (tab: "schema" | "shards" | "logs") => void;
  clearError: () => void;
}

// ─── Store ───────────────────────────────────────────────────────────────────

export const useProjectStore = create<ProjectState>((set, get) => ({
  // ── Initial state
  projects: [],
  activeProject: null,
  shards: [],
  logs: [],

  loadingProjects: false,
  loadingShards: false,
  loadingLogs: false,
  error: null,

  drawerOpen: false,
  activeTab: "schema",

  // ── Projects ─────────────────────────────────────────────────────────────

  fetchProjects: async () => {
    set({ loadingProjects: true, error: null });
    try {
      const res = await fetch("/api/projects");
      if (!res.ok) throw new Error("Failed to fetch projects");
      const data: Project[] = await res.json();
      set({ projects: data });
    } catch (e) {
      set({ error: (e as Error).message });
    } finally {
      set({ loadingProjects: false });
    }
  },

  createProject: async (name, description) => {
    set({ error: null });
    try {
      const res = await fetch("/api/projects", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, description }),
      });
      if (!res.ok) throw new Error("Failed to create project");
      const project: Project = await res.json();
      set((s) => ({ projects: [project, ...s.projects] }));
      return project;
    } catch (e) {
      set({ error: (e as Error).message });
      return null;
    }
  },

  setActiveProject: (project) => {
    set({ activeProject: project, shards: [], logs: [], activeTab: "schema" });
  },

  updateProjectSchema: async (projectId, schemaSql) => {
    try {
      const res = await fetch(`/api/projects/${projectId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ schema_sql: schemaSql }),
      });
      if (!res.ok) throw new Error("Failed to save schema");
      const updated: Project = await res.json();
      set((s) => ({
        projects: s.projects.map((p) => (p.id === projectId ? updated : p)),
        activeProject: s.activeProject?.id === projectId ? updated : s.activeProject,
      }));
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      return false;
    }
  },

  // ── Shards ───────────────────────────────────────────────────────────────

  fetchShards: async (projectId) => {
    set({ loadingShards: true });
    try {
      const res = await fetch(`/api/shards?project_id=${projectId}`);
      if (!res.ok) throw new Error("Failed to fetch shards");
      const data: Shard[] = await res.json();
      set({ shards: data });
    } catch (e) {
      set({ error: (e as Error).message });
    } finally {
      set({ loadingShards: false });
    }
  },

  deployShards: async (projectId, count) => {
    try {
      const res = await fetch("/api/shards/batch", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ project_id: projectId, count }),
      });
      if (!res.ok) throw new Error("Failed to deploy shards");
      // Re-fetch shards + update project shard_count
      await get().fetchShards(projectId);
      await get().fetchProjects();
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      return false;
    }
  },

  executeShard: async (shardId) => {
    try {
      const res = await fetch(`/api/shards/${shardId}/execute`, { method: "POST" });
      if (!res.ok) throw new Error("Execution failed");
      // Refresh shard list
      const { activeProject } = get();
      if (activeProject) await get().fetchShards(activeProject.id);
      return true;
    } catch (e) {
      set({ error: (e as Error).message });
      return false;
    }
  },

  // ── Logs ─────────────────────────────────────────────────────────────────

  fetchLogs: async (projectId) => {
    set({ loadingLogs: true });
    try {
      const res = await fetch(`/api/logs?project_id=${projectId}`);
      if (!res.ok) throw new Error("Failed to fetch logs");
      const data: Log[] = await res.json();
      set({ logs: data });
    } finally {
      set({ loadingLogs: false });
    }
  },

  // ── UI ───────────────────────────────────────────────────────────────────

  openDrawer: () => set({ drawerOpen: true }),
  closeDrawer: () => set({ drawerOpen: false }),
  setActiveTab: (tab) => set({ activeTab: tab }),
  clearError: () => set({ error: null }),
}));
