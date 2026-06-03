import { Pool } from "pg";

// Module-level singleton so the pool is reused across hot-reloads in dev
declare global {
  // eslint-disable-next-line no-var
  var __pgPool: Pool | undefined;
}

function createPool() {
  if (!process.env.DATABASE_URL) {
    throw new Error("DATABASE_URL is not set. Add it to .env.local");
  }
  return new Pool({ connectionString: process.env.DATABASE_URL });
}

const pool: Pool = global.__pgPool ?? (global.__pgPool = createPool());

export default pool;

/** Thin tagged-template helper so existing `sql\`...\`` call-sites work. */
export async function sql(
  strings: TemplateStringsArray,
  ...values: unknown[]
): Promise<Record<string, unknown>[]> {
  // Build a $1, $2, … parameterised query from the template literal
  let text = "";
  for (let i = 0; i < strings.length; i++) {
    text += strings[i];
    if (i < values.length) text += `$${i + 1}`;
  }
  const { rows } = await pool.query(text, values);
  return rows;
}
