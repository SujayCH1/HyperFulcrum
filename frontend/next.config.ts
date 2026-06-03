import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // keep server components from bundling pg (Node.js only)
  serverExternalPackages: ["pg"],
};

export default nextConfig;
