"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";
import { toast } from "sonner";
import {
  LayoutDashboard, Activity, Code2, Play, Database,
  Terminal as TermIcon, ArrowLeft, CheckCircle2, Circle,
  Server, Cpu, Zap, RefreshCw, Save, AlertCircle,
  ChevronDown, ChevronUp, Layers, ArrowDown, Check,
  RotateCcw,
} from "lucide-react";
import { useProjectStore } from "@/store/projectStore";
import type { Project, Shard, Log } from "@/types";

/* ─── Tux SVG ─────────────────────────────────────────────────────────── */
function TuxIcon({ size = 20 }: { size?: number }) {
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

const DEFAULT_SCHEMA = `-- Define your shard schema below
CREATE TABLE users (
  id         BIGSERIAL PRIMARY KEY,
  user_id    UUID NOT NULL DEFAULT gen_random_uuid(),
  email      VARCHAR(255) UNIQUE NOT NULL,
  username   VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_user_id ON users(user_id);
CREATE INDEX idx_users_email   ON users(email);

CREATE TABLE user_sessions (
  id         BIGSERIAL PRIMARY KEY,
  user_id    UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  token      VARCHAR(512) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);`;

/* ─── Walkthrough Overlay ─────────────────────────────────────────────── */
const WT_STEPS = [
  { n: "01", title: "Define Schema",    desc: "Write raw SQL CREATE TABLE statements — identical schema runs on every shard node.", color: "#10b981", border: "rgba(16,185,129,0.3)", bg: "rgba(16,185,129,0.08)", Icon: Code2 },
  { n: "02", title: "Configure Shards", desc: "Set how many shard nodes to provision. Each gets a unique key and host assignment.", color: "#3b82f6", border: "rgba(59,130,246,0.3)",  bg: "rgba(59,130,246,0.08)",  Icon: Server },
  { n: "03", title: "Execute & Monitor",desc: "Push schema to individual nodes, track health metrics and load in real time.",      color: "#8b5cf6", border: "rgba(139,92,246,0.3)", bg: "rgba(139,92,246,0.08)", Icon: Activity },
];

function WalkthroughOverlay({ onFinish }: { onFinish: () => void }) {
  const [visible, setVisible] = useState(0);
  const [showCTA, setShowCTA] = useState(false);

  useEffect(() => {
    const t = [
      setTimeout(() => setVisible(1), 350),
      setTimeout(() => setVisible(2), 1400),
      setTimeout(() => setVisible(3), 2450),
      setTimeout(() => setShowCTA(true), 3400),
    ];
    return () => t.forEach(clearTimeout);
  }, []);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-6"
      style={{ background: "rgba(3,3,3,0.96)", backdropFilter: "blur(14px)" }}>
      <div className="w-full max-w-[420px]">
        <div className="text-center mb-9 animate-fade-up">
          <div className="inline-flex items-center gap-2 text-[11px] font-mono rounded-full px-4 py-1.5 mb-4"
            style={{ color: "#a1a1aa", background: "#18181b", border: "1px solid #27272a" }}>
            <Zap size={10} style={{ color: "#10b981" }} />
            First-time Setup Walkthrough
          </div>
          <h2 className="text-[1.6rem] font-bold text-white mb-1.5">How it works</h2>
          <p className="text-sm" style={{ color: "#52525b" }}>Three steps to a distributed SQL cluster</p>
        </div>

        <div className="flex flex-col items-center">
          {WT_STEPS.map((s, i) => {
            const show = visible > i;
            return (
              <div key={s.n} className="flex flex-col items-center w-full">
                <div style={{
                  width: "100%", opacity: show ? 1 : 0,
                  transform: show ? "translateY(0) scale(1)" : "translateY(14px) scale(0.97)",
                  transition: "opacity 0.38s ease-out, transform 0.38s ease-out",
                  background: "#18181b", border: `1px solid ${show ? s.border : "#27272a"}`,
                  borderRadius: 14, padding: "16px 18px", display: "flex", alignItems: "flex-start", gap: 14,
                }}>
                  <div style={{ width: 40, height: 40, borderRadius: 10, flexShrink: 0, background: s.bg, border: `1px solid ${s.border}`, display: "flex", alignItems: "center", justifyContent: "center" }}>
                    <s.Icon size={17} style={{ color: s.color }} />
                  </div>
                  <div style={{ flex: 1 }}>
                    <div className="font-mono text-[9px] tracking-widest mb-1" style={{ color: "#52525b" }}>STEP {s.n}</div>
                    <div className="font-semibold text-sm text-white mb-1">{s.title}</div>
                    <p className="text-xs leading-relaxed" style={{ color: "#71717a" }}>{s.desc}</p>
                  </div>
                </div>
                {i < WT_STEPS.length - 1 && (
                  <div style={{ display: "flex", flexDirection: "column", alignItems: "center", marginTop: 4, marginBottom: 4, opacity: show ? 1 : 0, transition: "opacity 0.3s ease-out 0.25s" }}>
                    <div style={{ width: 1, height: 20, borderLeft: "2px dashed #27272a" }} />
                    <ArrowDown size={12} style={{ color: "#3f3f46", marginTop: -2 }} />
                  </div>
                )}
              </div>
            );
          })}
        </div>

        <div className="mt-7 flex flex-col items-center gap-3"
          style={{ opacity: showCTA ? 1 : 0, transform: showCTA ? "translateY(0)" : "translateY(10px)", transition: "opacity 0.35s, transform 0.35s" }}>
          <button onClick={onFinish} className="w-full font-bold py-3.5 rounded-xl font-mono text-sm text-white hover:bg-emerald-400 transition-all"
            style={{ background: "#059669", boxShadow: "0 4px 20px rgba(16,185,129,0.25)" }}>
            Begin Setup →
          </button>
          <button onClick={onFinish} className="text-xs font-mono hover:text-zinc-400 transition-colors" style={{ color: "#52525b" }}>
            Skip walkthrough
          </button>
        </div>
      </div>
    </div>
  );
}

/* ─── Project Wizard ─────────────────────────────────────────────────── */
function ProjectWizard({ project, onComplete }: { project: Project; onComplete: () => void }) {
  const [step, setStep]           = useState(1);
  const [schemaSql, setSchemaSql] = useState(DEFAULT_SCHEMA);
  const [shardCount, setShardCount] = useState(3);
  const [saving, setSaving]       = useState(false);
  const [deploying, setDeploying] = useState(false);

  const { updateProjectSchema, deployShards } = useProjectStore();

  const handleSaveSchema = async () => {
    setSaving(true);
    const ok = await updateProjectSchema(project.id, schemaSql);
    setSaving(false);
    if (ok) setStep(2);
    else toast.error("Failed to save schema");
  };

  const handleDeploy = async () => {
    setDeploying(true);
    const ok = await deployShards(project.id, shardCount);
    setDeploying(false);
    if (ok) { toast.success(`${shardCount} shard nodes deployed`); onComplete(); }
    else toast.error("Deployment failed");
  };

  const stepDone = step > 1;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-6"
      style={{ background: "rgba(3,3,3,0.96)", backdropFilter: "blur(14px)" }}>
      <div className="w-full max-w-2xl">
        {/* Step indicator */}
        <div className="flex items-center justify-center mb-8">
          {[{ n: 1, label: "Schema" }, { n: 2, label: "Shards" }].map((s, idx) => {
            const done   = step > s.n;
            const active = step === s.n;
            return (
              <React.Fragment key={s.n}>
                <div className="flex items-center gap-2">
                  <div className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold font-mono transition-all"
                    style={{ background: done || active ? "#10b981" : "transparent", border: done || active ? "1px solid #10b981" : "1px solid #3f3f46", color: done || active ? "#000" : "#52525b" }}>
                    {done ? <Check size={13} /> : s.n}
                  </div>
                  <span className="text-xs font-mono" style={{ color: active || done ? "#d4d4d8" : "#52525b" }}>{s.label}</span>
                </div>
                {idx === 0 && <div className="mx-3 h-px duration-500" style={{ width: 48, background: stepDone ? "#10b981" : "#27272a" }} />}
              </React.Fragment>
            );
          })}
        </div>

        {/* Step 1: Schema */}
        {step === 1 && (
          <div className="animate-fade-up">
            <div className="rounded-2xl overflow-hidden" style={{ background: "#111", border: "1px solid #27272a" }}>
              <div className="px-6 py-4 flex items-center justify-between" style={{ borderBottom: "1px solid #27272a", background: "#0d0d0d" }}>
                <div className="flex items-center gap-2.5">
                  <Code2 size={14} style={{ color: "#10b981" }} />
                  <span className="text-zinc-300 font-mono text-sm font-semibold">Schema Editor</span>
                </div>
                <span className="text-[10px] font-mono text-zinc-600">{project.name}</span>
              </div>
              <div className="relative">
                <textarea
                  value={schemaSql}
                  onChange={(e) => setSchemaSql(e.target.value)}
                  rows={16}
                  className="w-full bg-transparent px-6 py-5 text-emerald-300 font-mono text-[13px] leading-relaxed resize-none"
                  style={{ fontFamily: "'JetBrains Mono', 'Fira Code', monospace" }}
                  spellCheck={false}
                />
              </div>
            </div>
            <div className="mt-4 flex gap-3">
              <button
                onClick={handleSaveSchema}
                disabled={saving}
                className="flex-1 flex items-center justify-center gap-2 py-3.5 rounded-xl font-mono text-sm font-bold text-white transition-all"
                style={{ background: saving ? "#065f46" : "#059669", boxShadow: "0 4px 20px rgba(16,185,129,0.2)" }}
              >
                {saving ? <><span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />Saving...</> : <><Save size={14} />Save & Continue</>}
              </button>
              <button onClick={onComplete} className="px-6 py-3.5 rounded-xl font-mono text-sm text-zinc-500 border border-zinc-800 hover:border-zinc-700 hover:text-zinc-300 transition-all bg-zinc-900 hover:bg-zinc-800">
                Skip
              </button>
            </div>
          </div>
        )}

        {/* Step 2: Shards */}
        {step === 2 && (
          <div className="animate-fade-up">
            <div className="rounded-2xl p-8" style={{ background: "#111", border: "1px solid #27272a" }}>
              <div className="flex items-center gap-3 mb-6">
                <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ background: "rgba(59,130,246,0.1)", border: "1px solid rgba(59,130,246,0.25)" }}>
                  <Server size={18} className="text-blue-400" />
                </div>
                <div>
                  <div className="text-white font-semibold">Configure Shard Nodes</div>
                  <div className="text-zinc-600 text-xs font-mono">Set the number of nodes to provision</div>
                </div>
              </div>

              <div className="mb-8">
                <div className="flex items-center justify-between mb-3">
                  <label className="text-[10px] font-mono text-zinc-500 uppercase tracking-[0.18em]">Shard Count</label>
                  <span className="text-2xl font-bold text-white font-mono">{shardCount}</span>
                </div>
                <input
                  type="range" min={1} max={16} value={shardCount}
                  onChange={(e) => setShardCount(Number(e.target.value))}
                  className="w-full accent-blue-500"
                />
                <div className="flex justify-between text-[10px] font-mono text-zinc-700 mt-1">
                  <span>1</span><span>8</span><span>16</span>
                </div>
              </div>

              {/* Shard preview grid */}
              <div className="grid gap-2" style={{ gridTemplateColumns: `repeat(${Math.min(shardCount, 4)}, 1fr)` }}>
                {Array.from({ length: shardCount }, (_, i) => (
                  <div key={i} className="bg-zinc-900 border border-zinc-800 rounded-lg p-3 text-center">
                    <div className="w-2 h-2 rounded-full bg-blue-500 mx-auto mb-2" />
                    <div className="text-[9px] font-mono text-zinc-500">shard_{String.fromCharCode(97 + i)}</div>
                    <div className="text-[8px] font-mono text-zinc-700">node-{String(i + 1).padStart(2, "0")}</div>
                  </div>
                ))}
              </div>
            </div>

            <div className="mt-4 flex gap-3">
              <button onClick={() => setStep(1)} className="px-6 py-3.5 rounded-xl font-mono text-sm text-zinc-500 border border-zinc-800 hover:border-zinc-700 hover:text-zinc-300 transition-all bg-zinc-900">
                Back
              </button>
              <button
                onClick={handleDeploy}
                disabled={deploying}
                className="flex-1 flex items-center justify-center gap-2 py-3.5 rounded-xl font-mono text-sm font-bold text-white transition-all"
                style={{ background: deploying ? "#1d4ed8" : "#2563eb", boxShadow: "0 4px 20px rgba(59,130,246,0.2)" }}
              >
                {deploying
                  ? <><span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />Deploying...</>
                  : <><Layers size={14} />Deploy {shardCount} Shards</>}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/* ─── Stat Card ──────────────────────────────────────────────────────── */
function StatCard({ label, value, sub, color = "#10b981" }: { label: string; value: string | number; sub?: string; color?: string }) {
  return (
    <div className="bg-zinc-900/60 border border-zinc-800 rounded-xl p-5">
      <div className="text-[9px] font-mono text-zinc-600 uppercase tracking-[0.2em] mb-2">{label}</div>
      <div className="text-2xl font-bold font-mono" style={{ color }}>{value}</div>
      {sub && <div className="text-[10px] font-mono text-zinc-600 mt-1">{sub}</div>}
    </div>
  );
}

/* ─── Shard Row ──────────────────────────────────────────────────────── */
function ShardRow({ shard, onExecute }: { shard: Shard; onExecute: (id: number) => void }) {
  const [busy, setBusy]       = useState(false);
  const [expanded, setExpanded] = useState(false);

  const handleExec = async () => {
    setBusy(true);
    await onExecute(shard.id);
    setBusy(false);
  };

  const statusColor = shard.status === "active"   ? "#34d399"
                    : shard.status === "error"    ? "#f87171"
                    : "#71717a";

  return (
    <div className="border border-zinc-800 rounded-xl overflow-hidden" style={{ background: "rgba(9,9,9,0.6)" }}>
      <div className="flex items-center gap-4 px-5 py-4">
        <div className="flex-shrink-0">
          <div className="w-2 h-2 rounded-full" style={{ background: statusColor }} />
        </div>
        <TuxIcon size={16} />
        <div className="flex-1 min-w-0">
          <div className="text-zinc-200 font-mono text-sm font-semibold">{shard.shard_key}</div>
          <div className="text-zinc-600 text-[10px] font-mono">{shard.node_host}</div>
        </div>
        <div className="hidden sm:flex items-center gap-4 text-[10px] font-mono text-zinc-600">
          <span>load: {shard.metrics_load}%</span>
          {shard.schema_applied && <span className="text-emerald-500 flex items-center gap-1"><CheckCircle2 size={10} />schema</span>}
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleExec}
            disabled={busy || shard.schema_applied}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-mono text-[11px] font-semibold transition-all"
            style={{
              background: shard.schema_applied ? "transparent" : "rgba(16,185,129,0.1)",
              border: shard.schema_applied ? "1px solid #27272a" : "1px solid rgba(16,185,129,0.25)",
              color: shard.schema_applied ? "#3f3f46" : "#34d399",
              cursor: shard.schema_applied ? "not-allowed" : "pointer",
            }}
          >
            {busy ? <span className="w-3 h-3 border border-emerald-500/30 border-t-emerald-500 rounded-full animate-spin" /> : <Play size={10} />}
            {shard.schema_applied ? "Applied" : "Execute"}
          </button>
          <button onClick={() => setExpanded((x) => !x)} className="w-7 h-7 rounded-lg bg-zinc-900 border border-zinc-800 flex items-center justify-center text-zinc-500 hover:text-zinc-300 transition-all">
            {expanded ? <ChevronUp size={12} /> : <ChevronDown size={12} />}
          </button>
        </div>
      </div>

      {expanded && (
        <div className="px-5 pb-4 pt-0 border-t border-zinc-900 text-[11px] font-mono text-zinc-600 space-y-1">
          <div>ID: {shard.id}</div>
          <div>Created: {new Date(shard.created_at).toLocaleString()}</div>
          {shard.last_executed_at && <div>Last exec: {new Date(shard.last_executed_at).toLocaleString()}</div>}
        </div>
      )}
    </div>
  );
}

/* ─── Log Line ───────────────────────────────────────────────────────── */
function LogLine({ log }: { log: Log }) {
  const colors: Record<string, string> = {
    info:  "#a1a1aa",
    warn:  "#fbbf24",
    error: "#f87171",
    cmd:   "#34d399",
  };
  const prefixes: Record<string, string> = { info: "[INFO]", warn: "[WARN]", error: "[ERR ]", cmd: "[$]   " };
  return (
    <div className="flex items-start gap-3 py-1 font-mono text-[11px]">
      <span className="flex-shrink-0 text-zinc-700">{new Date(log.timestamp).toLocaleTimeString()}</span>
      <span className="flex-shrink-0" style={{ color: colors[log.level] ?? "#a1a1aa" }}>{prefixes[log.level] ?? "[LOG]"}</span>
      <span className="text-zinc-400 leading-relaxed">{log.message}</span>
    </div>
  );
}

/* ─── Dashboard Page ─────────────────────────────────────────────────── */
export default function DashboardPage() {
  const searchParams = typeof window !== "undefined" ? new URLSearchParams(window.location.search) : null;
  const projectId   = searchParams?.get("projectId");
  const firstTime   = searchParams?.get("firstTime") === "true";

  const {
    projects, activeProject, shards, logs,
    loadingShards, loadingLogs, error,
    fetchProjects, setActiveProject,
    fetchShards, fetchLogs, executeShard,
    activeTab, setActiveTab,
  } = useProjectStore();

  const [showWalkthrough, setShowWalkthrough] = useState(firstTime);
  const [showWizard, setShowWizard]           = useState(false);
  const [mounted, setMounted]                 = useState(false);
  const logsEndRef                            = useRef<HTMLDivElement>(null);

  useEffect(() => { setMounted(true); }, []);

  useEffect(() => {
    if (!mounted) return;
    fetchProjects();
  }, [mounted, fetchProjects]);

  useEffect(() => {
    if (!mounted || !projects.length) return;
    const target = projects.find((p) => String(p.id) === projectId) ?? projects[0];
    if (target) setActiveProject(target);
  }, [mounted, projects, projectId, setActiveProject]);

  useEffect(() => {
    if (!activeProject) return;
    fetchShards(activeProject.id);
    fetchLogs(activeProject.id);
  }, [activeProject, fetchShards, fetchLogs]);

  useEffect(() => {
    if (activeTab === "logs") logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs, activeTab]);

  const handleExecuteShard = useCallback(async (shardId: number) => {
    const ok = await executeShard(shardId);
    if (ok) { toast.success("Schema applied to shard"); fetchLogs(activeProject!.id); }
    else toast.error("Execution failed");
  }, [executeShard, activeProject, fetchLogs]);

  const handleWizardComplete = useCallback(() => {
    setShowWizard(false);
    if (activeProject) {
      fetchShards(activeProject.id);
      fetchLogs(activeProject.id);
    }
  }, [activeProject, fetchShards, fetchLogs]);

  if (!mounted) return null;

  const activeShards  = shards.filter((s) => s.status === "active").length;
  const appliedShards = shards.filter((s) => s.schema_applied).length;

  const tabs = [
    { id: "schema" as const,  label: "Schema",  Icon: Code2 },
    { id: "shards" as const,  label: "Shards",  Icon: Server },
    { id: "logs"   as const,  label: "Logs",    Icon: TermIcon },
  ];

  return (
    <div className="min-h-screen bg-[#050505] flex flex-col" style={{ fontFamily: "system-ui, sans-serif" }}>
      {/* Overlays */}
      {showWalkthrough && !showWizard && (
        <WalkthroughOverlay onFinish={() => { setShowWalkthrough(false); setShowWizard(true); }} />
      )}
      {showWizard && activeProject && (
        <ProjectWizard project={activeProject} onComplete={handleWizardComplete} />
      )}

      {/* Top nav */}
      <header className="flex items-center justify-between px-6 py-4 border-b border-zinc-900 flex-shrink-0">
        <div className="flex items-center gap-4">
          <a href="/" className="flex items-center gap-2 text-zinc-500 hover:text-zinc-300 transition-colors font-mono text-sm">
            <ArrowLeft size={14} /><span>Home</span>
          </a>
          <div className="w-px h-4 bg-zinc-800" />
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 bg-emerald-500 rounded flex items-center justify-center" style={{ boxShadow: "0 0 12px rgba(16,185,129,0.3)" }}>
              <Layers size={11} className="text-black" />
            </div>
            <span className="font-mono text-sm text-zinc-300 font-semibold">HyperFulcrum</span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          {error && (
            <div className="flex items-center gap-1.5 text-red-400 font-mono text-[11px]">
              <AlertCircle size={12} />{error}
            </div>
          )}
          <button
            onClick={() => activeProject && fetchShards(activeProject.id)}
            className="w-8 h-8 rounded-lg bg-zinc-900 border border-zinc-800 flex items-center justify-center text-zinc-500 hover:text-zinc-200 hover:border-zinc-700 transition-all"
          >
            <RefreshCw size={13} />
          </button>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar: project list */}
        <aside className="w-60 flex-shrink-0 border-r border-zinc-900 flex flex-col bg-[#070707]">
          <div className="px-4 py-3 border-b border-zinc-900">
            <p className="text-[9px] font-mono text-zinc-700 uppercase tracking-[0.2em]">Projects</p>
          </div>
          <div className="flex-1 overflow-y-auto py-2">
            {projects.map((p) => {
              const isActive = p.id === activeProject?.id;
              return (
                <button
                  key={p.id}
                  onClick={() => setActiveProject(p)}
                  className="w-full text-left px-4 py-3 flex items-center gap-3 transition-all group"
                  style={{ background: isActive ? "rgba(16,185,129,0.06)" : "transparent", borderLeft: isActive ? "2px solid #10b981" : "2px solid transparent" }}
                >
                  <Database size={13} className={isActive ? "text-emerald-400" : "text-zinc-600"} />
                  <div className="min-w-0 flex-1">
                    <div className={`font-mono text-xs truncate font-semibold ${isActive ? "text-zinc-100" : "text-zinc-400 group-hover:text-zinc-300"}`}>{p.name}</div>
                    <div className="text-zinc-700 text-[10px] font-mono">{p.shard_count ?? 0} shards</div>
                  </div>
                </button>
              );
            })}
          </div>
        </aside>

        {/* Main content */}
        <main className="flex-1 flex flex-col overflow-hidden">
          {activeProject ? (
            <>
              {/* Project header */}
              <div className="px-7 py-5 border-b border-zinc-900 flex items-center justify-between flex-shrink-0">
                <div>
                  <div className="flex items-center gap-3 mb-1">
                    <h1 className="text-lg font-bold text-white font-mono">{activeProject.name}</h1>
                    <span className="text-[9px] font-mono px-2 py-0.5 rounded"
                      style={{
                        color: activeProject.status === "active" ? "#34d399" : "#71717a",
                        background: activeProject.status === "active" ? "rgba(16,185,129,0.08)" : "#18181b",
                        border: `1px solid ${activeProject.status === "active" ? "rgba(16,185,129,0.2)" : "#27272a"}`,
                      }}>
                      {activeProject.status}
                    </span>
                  </div>
                  {activeProject.description && (
                    <p className="text-zinc-600 text-xs font-mono">{activeProject.description}</p>
                  )}
                </div>
                <button
                  onClick={() => setShowWizard(true)}
                  className="flex items-center gap-2 px-4 py-2 rounded-lg font-mono text-xs font-semibold transition-all text-emerald-400 border border-emerald-900/50 hover:bg-emerald-900/20"
                >
                  <RotateCcw size={12} />Re-run Wizard
                </button>
              </div>

              {/* Stats row */}
              <div className="px-7 py-4 grid grid-cols-4 gap-4 border-b border-zinc-900 flex-shrink-0">
                <StatCard label="Shards"    value={shards.length}           sub="provisioned"     />
                <StatCard label="Active"    value={activeShards}            sub="nodes online"    color="#34d399" />
                <StatCard label="Applied"   value={appliedShards}           sub="schema pushed"   color="#3b82f6" />
                <StatCard label="Logs"      value={logs.length}             sub="entries"         color="#a78bfa" />
              </div>

              {/* Tabs */}
              <div className="flex items-center gap-1 px-7 py-3 border-b border-zinc-900 flex-shrink-0">
                {tabs.map(({ id, label, Icon }) => {
                  const active = activeTab === id;
                  return (
                    <button
                      key={id}
                      onClick={() => setActiveTab(id)}
                      className="flex items-center gap-2 px-4 py-2 rounded-lg font-mono text-xs transition-all"
                      style={{
                        background: active ? "rgba(16,185,129,0.08)" : "transparent",
                        border: active ? "1px solid rgba(16,185,129,0.2)" : "1px solid transparent",
                        color: active ? "#34d399" : "#52525b",
                      }}
                    >
                      <Icon size={12} />
                      {label}
                    </button>
                  );
                })}
              </div>

              {/* Tab panels */}
              <div className="flex-1 overflow-y-auto p-7">
                {/* Schema tab */}
                {activeTab === "schema" && (
                  <div className="animate-fade-up">
                    {activeProject.schema_sql ? (
                      <div className="rounded-xl overflow-hidden" style={{ background: "#0d0d0d", border: "1px solid #27272a" }}>
                        <div className="px-5 py-3 flex items-center justify-between border-b border-zinc-900">
                          <div className="flex items-center gap-2 text-zinc-400 font-mono text-xs">
                            <Code2 size={12} /><span>schema.sql</span>
                          </div>
                          <button onClick={() => setShowWizard(true)} className="flex items-center gap-1.5 text-[10px] font-mono text-zinc-600 hover:text-zinc-300 transition-colors">
                            <Save size={10} />Edit
                          </button>
                        </div>
                        <pre className="px-6 py-5 text-emerald-300 font-mono text-[12px] leading-relaxed overflow-x-auto whitespace-pre-wrap">
                          {activeProject.schema_sql}
                        </pre>
                      </div>
                    ) : (
                      <div className="flex flex-col items-center justify-center py-20 text-center">
                        <Code2 size={32} className="text-zinc-800 mb-4" />
                        <p className="text-zinc-600 font-mono text-sm mb-4">No schema defined yet</p>
                        <button onClick={() => setShowWizard(true)} className="px-5 py-2.5 rounded-lg font-mono text-sm font-semibold text-white transition-all"
                          style={{ background: "#059669", boxShadow: "0 4px 16px rgba(16,185,129,0.2)" }}>
                          Open Wizard
                        </button>
                      </div>
                    )}
                  </div>
                )}

                {/* Shards tab */}
                {activeTab === "shards" && (
                  <div className="space-y-3 animate-fade-up">
                    {loadingShards ? (
                      <div className="flex items-center gap-2 text-zinc-600 font-mono text-sm py-8 justify-center">
                        <span className="w-4 h-4 border-2 border-zinc-700 border-t-zinc-400 rounded-full animate-spin" />Loading shards...
                      </div>
                    ) : shards.length === 0 ? (
                      <div className="flex flex-col items-center justify-center py-20 text-center">
                        <Server size={32} className="text-zinc-800 mb-4" />
                        <p className="text-zinc-600 font-mono text-sm mb-4">No shard nodes provisioned</p>
                        <button onClick={() => setShowWizard(true)} className="px-5 py-2.5 rounded-lg font-mono text-sm font-semibold text-white transition-all"
                          style={{ background: "#2563eb", boxShadow: "0 4px 16px rgba(59,130,246,0.2)" }}>
                          Deploy Shards
                        </button>
                      </div>
                    ) : (
                      shards.map((s) => (
                        <ShardRow key={s.id} shard={s} onExecute={handleExecuteShard} />
                      ))
                    )}
                  </div>
                )}

                {/* Logs tab */}
                {activeTab === "logs" && (
                  <div className="animate-fade-up rounded-xl overflow-hidden" style={{ background: "#080808", border: "1px solid #1a1a1a" }}>
                    <div className="px-5 py-3 border-b border-zinc-900 flex items-center justify-between">
                      <div className="flex items-center gap-2 text-zinc-500 font-mono text-xs">
                        <TermIcon size={12} /><span>orchestration.log</span>
                      </div>
                      <button onClick={() => fetchLogs(activeProject.id)} className="text-[10px] font-mono text-zinc-600 hover:text-zinc-300 transition-colors flex items-center gap-1">
                        <RefreshCw size={10} />Refresh
                      </button>
                    </div>
                    <div className="px-5 py-4 min-h-[300px] max-h-[500px] overflow-y-auto">
                      {loadingLogs ? (
                        <div className="text-zinc-700 font-mono text-xs py-4 text-center">Loading logs...</div>
                      ) : logs.length === 0 ? (
                        <div className="text-zinc-700 font-mono text-xs py-8 text-center">No logs yet</div>
                      ) : (
                        <>
                          {[...logs].reverse().map((l) => <LogLine key={l.id} log={l} />)}
                          <div ref={logsEndRef} />
                        </>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </>
          ) : (
            <div className="flex-1 flex flex-col items-center justify-center text-center">
              <Database size={40} className="text-zinc-800 mb-4" />
              <p className="text-zinc-600 font-mono text-sm">Select a project from the sidebar</p>
            </div>
          )}
        </main>
      </div>
    </div>
  );
}
