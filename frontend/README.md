# HyperFulcrum — Frontend

Container frontend for the HyperFulcrum SQL sharding orchestration backend.

## Stack
- **Next.js 15** (App Router)
- **Zustand** — all client state (replaces React Query / TanStack)
- **PostgreSQL** via `pg` pool — local Postgres, no cloud dependency
- **Tailwind CSS v3** + `sonner` for toasts
- **TypeScript strict**

## Setup

### 1. Prerequisites
- Node 18+
- PostgreSQL running locally (e.g. `brew services start postgresql`)

### 2. Install
```bash
npm install
```

### 3. Configure env
Copy `.env.local` and fill in your local Postgres URL:
```
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/hyperfrontend
```

### 4. Create the database
```bash
createdb hyperfrontend
```

### 5. Run migrations
```bash
npm run db:migrate
```

### 6. Start dev server
```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Project Structure
```
src/
  app/
    page.tsx            ← Home / project list
    dashboard/
      page.tsx          ← Main orchestration dashboard
    api/
      projects/         ← CRUD for projects
      shards/           ← Shard provisioning + execution
      logs/             ← Orchestration log read/write
  store/
    projectStore.ts     ← Zustand store (all state + fetchers)
  lib/
    db.ts               ← pg Pool singleton + sql helper
    dbInit.ts           ← CREATE TABLE IF NOT EXISTS on first request
  types/
    index.ts            ← Project, Shard, Log interfaces
```

## Zustand Store
`useProjectStore` exposes all data and actions. Consume it anywhere:

```ts
import { useProjectStore } from "@/store/projectStore";

const { projects, fetchProjects, createProject } = useProjectStore();
```

No providers needed — Zustand is provider-free.

## Adding More Features
The wizard currently covers schema definition + shard deployment.
Future steps (e.g. shard rebalancing, query routing config, health polling)
can be added as extra wizard steps or new dashboard tabs without touching
the store shape — just extend `projectStore.ts` with new state slices
and API calls.
