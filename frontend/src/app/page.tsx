"use client";

import React, { useState, useEffect } from "react";
import { toast } from "sonner";
import {
  Layers, Plus, Database, ChevronRight, Cpu, X,
  ArrowRight, GitBranch, Activity,
} from "lucide-react";
import { useProjectStore } from "@/store/projectStore";
import type { Project } from "@/types";

/* ─── Tux SVG ──────────────────────────────────────────────────────────── */
function TuxIcon({ size = 18 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 48 55" fill="none">
      <ellipse cx="24" cy="14" rx="13" ry="13" fill="#f5f5f5" />
      <ellipse cx="17.5" cy="12" rx="3.5" ry="4.5" fill="#222" />
      <ellipse cx="30.5" cy="12" rx="3.5" ry="4.5" fill="#222" />
      <ellipse cx="18" cy="11.5" rx="1.8" ry="2.2" fill="#fff" />
      <ellipse cx="31" cy="11.5" rx="1.8" ry="2.2" fill="#fff" />
      <ellipse cx="18.5" cy="12" rx="1" ry="1.2" fill="#111" />
      <ellipse cx="31.5" cy="12" rx="1" ry="1.2" fill="#111" />
      <path d="M20 20 Q24 24 28 20" stroke="#e8a020" strokeWidth="1.5" strokeLinecap="round" fill="none" />
      <ellipse cx="11" cy="15" rx="4.5" ry="3.5" fill="#f5f5f5" />
      <ellipse cx="37" cy="15" rx="4.5" ry="3.5" fill="#f5f5f5" />
      <rect x="11" y="25" width="26" height="22" rx="7" fill="#f0a030" />
      <ellipse cx="24" cy="26" rx="12" ry="5" fill="#f0a030" />
      <rect x="9" y="40" width="9" height="11" rx="4" fill="#2a2a2a" />
      <rect x="30" y="40" width="9" height="11" rx="4" fill="#2a2a2a" />
    </svg>
  );
}

/* ─── Create Drawer ─────────────────────────────────────────────────────── */
function CreateDrawer({ open, onClose, onSuccess }: {
  open: boolean;
  onClose: () => void;
  onSuccess: (p: Project) => void;
}) {
  const [name, setName]   = useState("");
  const [desc, setDesc]   = useState("");
  const [busy, setBusy]   = useState(false);
  const createProject     = useProjectStore((s) => s.createProject);

  useEffect(() => {
    document.body.style.overflow = open ? "hidden" : "";
    return () => { document.body.style.overflow = ""; };
  }, [open]);

  const handleSubmit = async () => {
    const trimmed = name.trim();
    if (!trimmed) { toast.error("Project name is required"); return; }
    setBusy(true);
    const project = await createProject(trimmed, desc.trim());
    setBusy(false);
    if (project) {
      toast.success("Project created");
      setName(""); setDesc("");
      onSuccess(project);
    } else {
      toast.error("Failed to create project");
    }
  };

  const handleKey = (e: React.KeyboardEvent) => {
    if ((e.metaKey || e.ctrlKey) && e.key === "Enter") handleSubmit();
    if (e.key === "Escape") onClose();
  };

  if (!open) return null;

  const hasName  = Boolean(name.trim());
  const glowStyle = hasName ? { boxShadow: "0 0 0 2px rgba(16,185,129,0.3)" } : {};

  return (
    <div className="fixed inset-0 z-50 flex" onKeyDown={handleKey}>
      {/* Backdrop */}
      <div
        className="absolute inset-0"
        style={{ background: "rgba(0,0,0,0.8)", backdropFilter: "blur(10px)" }}
        onClick={onClose}
      />

      {/* Drawer panel */}
      <div
        className="absolute right-0 top-0 bottom-0 flex flex-col bg-[#0b0b0b] border-l border-zinc-800 animate-slide-in"
        style={{ width: "min(880px, 100%)" }}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-7 py-4 border-b border-zinc-900 flex-shrink-0">
          <div className="flex items-center gap-3">
            <div className="w-7 h-7 rounded-lg flex items-center justify-center"
              style={{ background: "rgba(16,185,129,0.12)", border: "1px solid rgba(16,185,129,0.25)" }}>
              <Database size={13} className="text-emerald-400" />
            </div>
            <span className="text-zinc-200 font-mono text-sm font-semibold">Initialize Sharded Database</span>
          </div>
          <button
            onClick={onClose}
            className="w-7 h-7 rounded-lg bg-zinc-900 border border-zinc-800 flex items-center justify-center text-zinc-500 hover:text-zinc-200 hover:border-zinc-700 transition-all"
          >
            <X size={13} />
          </button>
        </div>

        {/* Body: two columns */}
        <div className="flex flex-1 overflow-hidden">
          {/* Left: form */}
          <div className="flex-1 flex flex-col px-8 py-7 overflow-y-auto border-r border-zinc-900/80">
            <div className="mb-7">
              <h2 className="text-[1.35rem] font-bold text-white mb-2 tracking-tight">New Project</h2>
              <p className="text-zinc-500 text-sm leading-relaxed max-w-sm">
                A project is a distributed SQL database split across shard nodes. Name it, define a schema, and deploy.
              </p>
            </div>

            <div className="space-y-5 flex-1">
              <div>
                <label className="block text-[10px] font-mono text-zinc-500 uppercase tracking-[0.18em] mb-2">
                  Project Name <span className="text-red-500/80">*</span>
                </label>
                <div className="relative">
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    autoFocus
                    placeholder="users-shard-cluster"
                    className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-3.5 text-white font-mono text-sm placeholder-zinc-700 transition-all"
                    style={glowStyle}
                  />
                  {hasName && <span className="absolute right-3.5 top-1/2 -translate-y-1/2 w-1.5 h-1.5 rounded-full bg-emerald-500" />}
                </div>
                <p className="text-zinc-700 text-[10px] font-mono mt-1.5">e.g. analytics-db, user-cluster, events-shard</p>
              </div>

              <div>
                <label className="block text-[10px] font-mono text-zinc-500 uppercase tracking-[0.18em] mb-2">
                  Description <span className="text-zinc-700 normal-case">(optional)</span>
                </label>
                <textarea
                  value={desc}
                  onChange={(e) => setDesc(e.target.value)}
                  placeholder="Horizontal sharding for user auth data across 4 regions..."
                  rows={3}
                  className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-3.5 text-white font-mono text-sm placeholder-zinc-700 resize-none transition-all"
                />
              </div>

              {/* Steps preview */}
              <div className="bg-zinc-900/50 border border-zinc-800/60 rounded-xl p-5">
                <p className="text-[9px] font-mono text-zinc-600 uppercase tracking-[0.2em] mb-3.5">What happens next</p>
                <div className="space-y-2.5">
                  {[
                    { n: "01", t: "Walkthrough",     d: "See how the system works" },
                    { n: "02", t: "Write Schema",    d: "Define SQL tables & indexes" },
                    { n: "03", t: "Configure Shards",d: "Set node count & deploy" },
                  ].map((s) => (
                    <div key={s.n} className="flex items-center gap-3">
                      <span className="text-[9px] font-mono text-zinc-700 w-4 flex-shrink-0">{s.n}</span>
                      <span className="text-zinc-400 text-xs font-medium">{s.t}</span>
                      <span className="text-zinc-700 text-xs">·</span>
                      <span className="text-zinc-600 text-xs">{s.d}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="mt-7 space-y-3">
              <button
                onClick={handleSubmit}
                disabled={busy || !hasName}
                className="w-full flex items-center justify-center gap-2.5 py-3.5 rounded-xl font-mono text-sm font-bold text-white transition-all"
                style={{
                  background: hasName && !busy ? "#059669" : "#065f46",
                  opacity: !hasName ? 0.4 : 1,
                  cursor: !hasName ? "not-allowed" : "pointer",
                  boxShadow: hasName ? "0 4px 20px rgba(16,185,129,0.2)" : "none",
                }}
              >
                {busy ? (
                  <><span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />Initializing...</>
                ) : (
                  <><Layers size={14} />Initialize Project<ArrowRight size={13} /></>
                )}
              </button>
              <button
                onClick={onClose}
                className="w-full py-3 rounded-xl text-zinc-500 hover:text-zinc-300 font-mono text-sm border border-zinc-800 hover:border-zinc-700 bg-zinc-900 hover:bg-zinc-800 transition-all"
              >
                Cancel
              </button>
              <p className="text-[10px] text-zinc-700 font-mono text-center">⌘ + Enter to submit</p>
            </div>
          </div>

          {/* Right: live preview */}
          <div className="w-[300px] flex-shrink-0 flex flex-col px-6 py-7 overflow-y-auto">
            <p className="text-[9px] font-mono text-zinc-600 uppercase tracking-[0.22em] mb-5">Live Preview</p>

            <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-4 mb-4">
              <div className="flex items-center gap-2.5 mb-3">
                <div className="w-6 h-6 rounded-md flex items-center justify-center flex-shrink-0"
                  style={{ background: "rgba(16,185,129,0.1)", border: "1px solid rgba(16,185,129,0.2)" }}>
                  <Database size={11} className="text-emerald-400" />
                </div>
                <span className="font-mono text-xs text-zinc-300 font-semibold truncate">{name || "project-name"}</span>
              </div>
              {desc && <p className="text-zinc-600 text-[10px] font-mono mb-3 leading-relaxed line-clamp-2">{desc}</p>}
              <div className="flex gap-1.5">
                <span className="text-[9px] font-mono text-zinc-500 bg-zinc-800 rounded px-2 py-0.5">idle</span>
                <span className="text-[9px] font-mono text-zinc-600 rounded px-2 py-0.5">0 shards</span>
              </div>
            </div>

            {/* Architecture mini-diagram */}
            <div className="bg-zinc-900/40 border border-zinc-900 rounded-xl p-4 mb-4">
              <p className="text-[9px] font-mono text-zinc-700 mb-4 uppercase tracking-widest">Architecture</p>
              <div className="bg-zinc-800/60 border border-zinc-700/60 rounded-lg px-3 py-2 text-center mb-2">
                <span className="text-[10px] font-mono text-zinc-400">SQL Router Layer</span>
              </div>
              <div className="flex justify-center my-2"><div className="w-px h-4 bg-zinc-700" /></div>
              <div className="grid grid-cols-3 gap-1.5">
                {["shard_a", "shard_b", "shard_c"].map((k) => (
                  <div key={k} className="bg-zinc-900 border border-zinc-800 rounded-md px-1.5 py-2 text-center">
                    <div className="w-1.5 h-1.5 rounded-full bg-zinc-700 mx-auto mb-1" />
                    <span className="text-[8px] font-mono text-zinc-600">{k}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-2.5">
              <div className="flex items-center gap-2.5"><TuxIcon size={16} /><span className="text-[10px] font-mono text-zinc-600">Linux PostgreSQL nodes</span></div>
              <div className="flex items-center gap-2.5"><GitBranch size={11} className="text-zinc-600" /><span className="text-[10px] font-mono text-zinc-600">Horizontal sharding</span></div>
              <div className="flex items-center gap-2.5"><Activity size={11} className="text-zinc-600" /><span className="text-[10px] font-mono text-zinc-600">Real-time health metrics</span></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

/* ─── Home Page ─────────────────────────────────────────────────────────── */
export default function HomePage() {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { projects, fetchProjects, loadingProjects } = useProjectStore();

  useEffect(() => { fetchProjects(); }, [fetchProjects]);

  const handleSuccess = (project: Project) => {
    setDrawerOpen(false);
    window.location.href = `/dashboard?projectId=${project.id}&firstTime=true`;
  };

  return (
    <div className="min-h-screen bg-[#050505] text-zinc-100 flex flex-col overflow-x-hidden">
      {/* Grid pattern */}
      <div
        className="fixed inset-0 pointer-events-none"
        style={{
          backgroundImage:
            "linear-gradient(rgba(16,185,129,0.024) 1px,transparent 1px)," +
            "linear-gradient(90deg,rgba(16,185,129,0.024) 1px,transparent 1px)",
          backgroundSize: "52px 52px",
        }}
      />

      {/* Header */}
      <header className="relative z-10 px-8 py-5 flex items-center justify-between border-b border-zinc-900/80 animate-fade-up">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-emerald-500 rounded-lg flex items-center justify-center"
            style={{ boxShadow: "0 0 20px rgba(16,185,129,0.35)" }}>
            <Layers size={14} className="text-black" />
          </div>
          <span className="font-mono text-sm text-zinc-200 font-bold tracking-[0.18em] uppercase">HyperFulcrum</span>
        </div>
        <div className="flex items-center gap-2 text-[11px] font-mono text-zinc-600">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-blink" />
          Orchestration Engine · v0.1
        </div>
      </header>

      <main className="relative z-10 flex-1 flex flex-col items-center justify-center px-6 py-16">
        {/* Badge */}
        <div className="inline-flex items-center gap-2 bg-zinc-900/90 border border-zinc-800 rounded-full px-4 py-1.5 mb-10 text-xs text-emerald-400 font-mono animate-fade-up">
          <Cpu size={10} />
          SQL Sharding Orchestration Engine
        </div>

        {/* Hero text */}
        <div className="text-center mb-14 max-w-2xl animate-fade-up">
          <h1 className="text-5xl font-extrabold text-white mb-5 leading-[1.07] tracking-tight">
            Distribute your database
            <br />
            <span className="text-emerald-400">across infinite nodes</span>
          </h1>
          <p className="text-zinc-500 text-lg leading-relaxed max-w-lg mx-auto">
            Define schema, set shard count, and deploy across Linux nodes in minutes.
          </p>
        </div>

        {/* CTA */}
        <div className="relative flex flex-col items-center animate-fade-up">
          {/* Curly arrow */}
          <div className="absolute pointer-events-none" style={{ top: "-86px", right: "-215px" }}>
            <svg width="200" height="108" viewBox="0 0 200 108" fill="none">
              <path
                d="M 183 16 C 163 7 140 18 122 34 C 104 50 93 70 75 85"
                stroke="#10b981" strokeWidth="1.8" strokeDasharray="7 5" strokeLinecap="round"
                className="animate-dash-flow"
              />
              <path d="M 64 89 L 78 79 L 76 93 Z" fill="#10b981" opacity="0.8" />
            </svg>
            <div className="absolute text-[11px] text-emerald-400 font-mono whitespace-nowrap px-2 py-0.5 rounded animate-fade-up"
              style={{ top: 0, right: 0, background: "rgba(5,5,5,0.85)", border: "1px solid rgba(16,185,129,0.2)" }}>
              ↗ start here
            </div>
          </div>

          <button
            onClick={() => setDrawerOpen(true)}
            className="group flex items-center gap-5 border-2 border-dashed rounded-2xl px-10 py-8 transition-all duration-300 hover:border-emerald-500"
            style={{ borderColor: "#3f3f46", background: "rgba(8,8,8,0.85)", backdropFilter: "blur(12px)" }}
          >
            <div className="w-14 h-14 rounded-xl flex items-center justify-center flex-shrink-0 transition-all"
              style={{ background: "rgba(16,185,129,0.08)", border: "1px solid rgba(16,185,129,0.22)" }}>
              <Plus size={28} className="text-emerald-400" />
            </div>
            <div className="text-left">
              <div className="text-white font-bold text-xl mb-1">Create New Project</div>
              <div className="text-zinc-500 text-sm font-mono">Initialize a horizontally sharded SQL database</div>
            </div>
            <ArrowRight size={20} className="text-zinc-600 group-hover:text-emerald-400 transition-colors ml-4" />
          </button>

          <div className="flex flex-wrap items-center justify-center gap-6 mt-7 animate-fade-up">
            {["Schema Management", "Auto Sharding", "Node Monitoring", "CLI Access"].map((f) => (
              <span key={f} className="text-xs text-zinc-600 font-mono flex items-center gap-1.5">
                <span className="w-1 h-1 rounded-full bg-emerald-800 flex-shrink-0" />
                {f}
              </span>
            ))}
          </div>
        </div>

        {/* Recent projects list */}
        {projects.length > 0 && (
          <div className="mt-20 w-full max-w-lg animate-fade-up">
            <p className="text-[9px] font-mono text-zinc-700 uppercase tracking-[0.24em] mb-4 text-center">Recent Projects</p>
            <div className="space-y-2">
              {projects.slice(0, 6).map((p) => {
                const isActive = p.status === "active";
                const badge = isActive
                  ? { color: "#34d399", bg: "rgba(16,185,129,0.08)", border: "rgba(16,185,129,0.2)" }
                  : { color: "#71717a", bg: "#18181b", border: "#27272a" };
                return (
                  <a
                    key={p.id}
                    href={`/dashboard?projectId=${p.id}`}
                    className="flex items-center justify-between border border-zinc-800 hover:border-zinc-700 rounded-xl px-5 py-3.5 transition-all group"
                    style={{ background: "rgba(9,9,9,0.65)" }}
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-7 h-7 rounded-lg bg-zinc-800 flex items-center justify-center flex-shrink-0">
                        <Database size={12} className="text-emerald-400" />
                      </div>
                      <div>
                        <div className="text-zinc-300 font-mono text-sm font-semibold">{p.name}</div>
                        <div className="text-zinc-600 text-[10px] font-mono">
                          {p.shard_count ?? 0} shards{p.description ? " · " + p.description.slice(0, 38) : ""}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 flex-shrink-0 ml-3">
                      <span className="text-[9px] font-mono px-2 py-0.5 rounded"
                        style={{ color: badge.color, background: badge.bg, border: `1px solid ${badge.border}` }}>
                        {p.status}
                      </span>
                      <ChevronRight size={13} className="text-zinc-600 group-hover:text-zinc-400 transition-colors" />
                    </div>
                  </a>
                );
              })}
            </div>
            {loadingProjects && (
              <p className="text-center text-zinc-700 font-mono text-[10px] mt-4">Fetching...</p>
            )}
          </div>
        )}
      </main>

      <CreateDrawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        onSuccess={handleSuccess}
      />
    </div>
  );
}
